package testing

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"

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
}

type MockOptions struct {
	Callback func(db *gorm.DB, g *gin.Engine)
	Models   []any
	BasePath string
}

func NewMock(options *MockOptions) Mock {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	dialect := sqlite.Open(":memory:")
	db, _ := gorm.Open(dialect, &gorm.Config{Logger: logger.New(log.New(os.Stdout, "\n", log.LstdFlags), logger.Config{Colorful: true, LogLevel: logger.LogLevel(logger.Info)})})
	// db, _ := gorm.Open(dialect)
	db.AutoMigrate(options.Models...)

	if options.Callback != nil {
		options.Callback(db, engine)
	}

	m := Mock{
		Engine:   engine,
		Response: *httptest.NewRecorder(),
		Database: db,
	}

	return m
}

type Response struct {
	Status int
	Body   string
	Json   map[string]any
}

func (m *Mock) request(method, endpoint string) Response {
	var response map[string]any
	m.Response = *httptest.NewRecorder()

	req, _ := http.NewRequest(method, endpoint, nil)
	m.Engine.ServeHTTP(&m.Response, req)

	json.Unmarshal(m.Response.Body.Bytes(), &response)

	return Response{Body: m.Response.Body.String(), Json: response, Status: m.Response.Code}
}

func (m *Mock) GET(endpoint string) Response {
	return m.request("GET", fmt.Sprintf("%s%s", m.BasePath, endpoint))
}

func (m *Mock) POST(endpoint string, form gin.H) Response {
	return m.request("POST", fmt.Sprintf("%s%s", m.BasePath, endpoint))
}
