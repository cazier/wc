package testing

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var here string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	here = path.Dir(filename)
	root := path.Join(here, "..")
	os.Chdir(root)
}

func Path(name string) string {
	return path.Join(here, name)
}

type Mock struct {
	Engine   *gin.Engine
	Database *gorm.DB
	Response httptest.ResponseRecorder
	BasePath string
	models   []any
	ctx      *gin.Context
}

type MockOptions struct {
	Callback func(db *gorm.DB, g *gin.Engine)
	Models   []any
	BasePath string
}

func NewMock(options *MockOptions) Mock {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	m := Mock{
		Engine: engine,
		models: make([]any, len(options.Models)),
	}
	copy(m.models, options.Models)

	m.Database = m.OpenDB()

	if options.Callback != nil {
		options.Callback(m.Database, m.Engine)
	}

	return m
}

func (m Mock) OpenDB() *gorm.DB {
	dialect := sqlite.Open(":memory:")

	if _, ok := os.LookupEnv("VERBOSE_TESTING"); ok {
		m.Database, _ = gorm.Open(
			dialect,
			&gorm.Config{
				Logger: logger.New(
					log.New(os.Stdout, "\n", log.LstdFlags),
					logger.Config{
						Colorful: true,
						LogLevel: logger.LogLevel(logger.Info),
					},
				),
			},
		)
	} else {
		m.Database, _ = gorm.Open(dialect)
	}
	m.Database.AutoMigrate(m.models...)

	return m.Database
}

func (m Mock) CloseDB() {
	sql, _ := m.Database.DB()

	sql.Close()
}

type Response struct {
	Status  int
	Body    string
	Headers http.Header
	Context *gin.Context
	Json    map[string]any
}

type RequestOptions struct {
	Cookies []map[string]string
	Form    map[string]any
}

func (m *Mock) request(method, endpoint string, options ...RequestOptions) Response {

	var response map[string]any
	m.Response = *httptest.NewRecorder()
	m.ctx = gin.CreateTestContextOnly(&m.Response, m.Engine)

	m.ctx.Request, _ = http.NewRequest(method, fmt.Sprintf("%s%s", m.BasePath, endpoint), nil)
	option := m.applyOptions(options)
	m.encodeForm(option)
	m.Engine.ServeHTTP(&m.Response, m.ctx.Request)

	json.Unmarshal(m.Response.Body.Bytes(), &response)

	return Response{Body: m.Response.Body.String(), Json: response, Headers: m.Response.Header(), Status: m.Response.Code, Context: m.ctx}
}

func (m *Mock) applyOptions(options []RequestOptions) RequestOptions {
	resp := RequestOptions{}

	if len(options) != 1 {
		return resp
	}

	option := options[0]

	for _, cookie := range option.Cookies {
		m.ctx.SetCookie(cookie["key"], cookie["value"], 0, "", "/", false, false)
	}

	return option
}

func (m *Mock) encodeForm(options RequestOptions) {
	if options.Form == nil {
		return
	}

	data := url.Values{}

	for key, value := range options.Form {
		data.Set(key, value.(string))
	}

	encoded := strings.NewReader(data.Encode())
	m.ctx.Request, _ = http.NewRequest(m.ctx.Request.Method, m.ctx.Request.URL.Path, encoded)
	m.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func (m *Mock) GET(endpoint string, options ...RequestOptions) Response {
	return m.request(http.MethodGet, endpoint, options...)
}

func (m *Mock) POST(endpoint string, options ...RequestOptions) Response {
	return m.request(http.MethodPost, endpoint, options...)
}
