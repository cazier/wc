package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
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

func (m *Mock) POST(endpoint string) Response {
	return m.request("POST", fmt.Sprintf("%s%s", m.BasePath, endpoint))
}
