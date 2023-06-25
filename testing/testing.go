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
	"regexp"
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
	Engine         *gin.Engine
	Database       *gorm.DB
	Response       httptest.ResponseRecorder
	BasePath       string
	models         []any
	ctx            *gin.Context
	cookieCallback func(m *Mock)
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
				TranslateError: true,
			},
		)
	} else {
		m.Database, _ = gorm.Open(dialect, &gorm.Config{TranslateError: true})
	}
	m.Database.AutoMigrate(m.models...)

	return m.Database
}

type Response struct {
	Status  int
	Body    string
	Headers http.Header
	Context *gin.Context
	Json    map[string]any
}

type RequestOptions struct {
	Cookies map[string]string
	Form    map[string]any
}

func (m *Mock) request(method, endpoint string, options ...RequestOptions) Response {
	var response map[string]any
	m.Response = *httptest.NewRecorder()
	m.ctx = gin.CreateTestContextOnly(&m.Response, m.Engine)

	m.ctx.Request, _ = http.NewRequest(method, fmt.Sprintf("%s%s", m.BasePath, endpoint), nil)

	if m.cookieCallback != nil {
		m.cookieCallback(m)
		m.cookieCallback = nil
	}

	for _, option := range options {
		encodeForm(m.ctx, option.Form)
	}

	m.Engine.ServeHTTP(&m.Response, m.ctx.Request)

	json.Unmarshal(m.Response.Body.Bytes(), &response)

	return Response{
		Body:    m.Response.Body.String(),
		Json:    response,
		Headers: m.Response.Header(),
		Status:  m.Response.Code,
	}
}

func (m *Mock) WithSetCookie(from Response) *Mock {
	setCookie := from.Headers.Get("Set-Cookie")

	pattern := regexp.MustCompile("session=(.*?);")
	cookie := pattern.FindStringSubmatch(setCookie)

	m.cookieCallback = func(m *Mock) {
		m.ctx.Request.AddCookie(&http.Cookie{Name: "session", Value: cookie[1]})
	}

	return m
}

func encodeForm(ctx *gin.Context, values map[string]any) {
	if values == nil {
		return
	}

	data := url.Values{}

	for key, value := range values {
		data.Set(key, value.(string))
	}

	encoded := strings.NewReader(data.Encode())
	if ctx.Request != nil {
		ctx.Request, _ = http.NewRequest(ctx.Request.Method, ctx.Request.URL.Path, encoded)
	} else {
		ctx.Request, _ = http.NewRequest(http.MethodPost, "", encoded)
	}
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func (m *Mock) GET(endpoint string, options ...RequestOptions) Response {
	return m.request(http.MethodGet, endpoint, options...)
}

func (m *Mock) POST(endpoint string, options ...RequestOptions) Response {
	return m.request(http.MethodPost, endpoint, options...)
}

func (m *Mock) Form(data map[string]any) *gin.Context {
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	encodeForm(context, data)

	return context
}
