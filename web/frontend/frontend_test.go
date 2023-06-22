package frontend

import (
	"net/http"
	"os"
	"testing"

	"github.com/cazier/wc/db/models"
	test "github.com/cazier/wc/testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var m test.Mock

func TestMain(tm *testing.M) {
	m = test.NewMock(
		&test.MockOptions{
			Callback: Init,
			Models: []any{
				&models.User{},
			},
		},
	)

	os.Exit(tm.Run())
}

func TestRegisterGet(t *testing.T) {
	assert := assert.New(t)

	response := m.GET("/register")

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "Register")
}

func TestLoginGet(t *testing.T) {
	assert := assert.New(t)

	response := m.GET("/login")

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "Login")
}

func TestRegisterPost(t *testing.T) {
	assert := assert.New(t)

	response := m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"email": "register@email.com", "password": "password", "confirm": "password"},
		},
	)

	assert.Equal(http.StatusFound, response.Status)
	assert.Equal("/secret", response.Headers.Get("Location"))
}

func TestRegisterBad(t *testing.T) {
	assert := assert.New(t)

	response := m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"email": "register@email.com", "confirm": "password"},
		},
	)
	assert.Equal(http.StatusNotAcceptable, response.Status)
	assert.Contains(response.Body, "invalid username or password")

	response = m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"email": "register@email.com", "password": "password", "confirm": "passw0rd"},
		},
	)
	assert.Equal(http.StatusNotAcceptable, response.Status)
	assert.Contains(response.Body, "passwords do not match")

	m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"email": "register@email.com", "password": "password", "confirm": "password"},
		},
	)
	response = m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"email": "register@email.com", "password": "password", "confirm": "password"},
		},
	)
	assert.Equal(http.StatusNotAcceptable, response.Status)
	assert.Contains(response.Body, "an account with this email address already exists")
}

func TestLoginPost(t *testing.T) {
	assert := assert.New(t)

	m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"email": "login@email.com", "password": "password", "confirm": "password"},
		},
	)

	response := m.POST(
		"/login",
		test.RequestOptions{
			Form: gin.H{"email": "login@email.com", "password": "password"},
		},
	)

	assert.Equal(http.StatusFound, response.Status)
	assert.Equal("/secret", response.Headers.Get("Location"))
}
