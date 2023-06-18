package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"

	database "github.com/cazier/wc/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	Response httptest.ResponseRecorder
	BasePath string
}

type MockOptions struct {
	BasePath      string
	BasePathGroup *gin.RouterGroup
	Callback      func(db *gorm.DB, g *gin.Engine)
}

func NewMock(options *MockOptions) Mock {
	gin.SetMode(gin.ReleaseMode)

	database.InitSqlite(&database.SqliteDBOptions{Memory: true, LogLevel: 3, Purge: database.No})

	engine := gin.New()
	db := database.Database

	options.Callback(db, engine)

	m := Mock{
		Engine:   engine,
		Response: *httptest.NewRecorder(),
	}

	if options.BasePathGroup != nil {
		m.BasePath = options.BasePathGroup.BasePath()
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
