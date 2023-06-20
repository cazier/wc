package frontend

import (
	"net/http"
	"testing"

	"github.com/cazier/wc/db/models"
	test "github.com/cazier/wc/testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var m test.Mock

func init() {
	m = test.NewMock(
		&test.MockOptions{
			Callback: Init,
			Models: []any{
				&models.User{},
			},
		},
	)
}

func TestLoginGet(t *testing.T) {
	assert := assert.New(t)

	response := m.GET("/login")

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "Login")
}

func TestRegisterGet(t *testing.T) {
	assert := assert.New(t)

	response := m.GET("/register")

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "Register")
}

func TestRegisterPost(t *testing.T) {
	assert := assert.New(t)

	response := m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{
				"email":    "register@email.com",
				"password": "password",
				"confirm":  "password",
			},
		},
	)

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "Register")
}

// func TestLoginPost(t *testing.T) {
// 	assert := assert.New(t)

// 	response := m.POST("/login", test.RequestOptions{Form: gin.H{"email": "post@email.com"}})
// 	assert.Equal(http.StatusNotAcceptable, response.Status)
// }
