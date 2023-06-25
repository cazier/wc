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
			Form: gin.H{"name": "register", "email": "register@email.com", "password": "pass", "confirm": "pass"},
		},
	)

	assert.Equal(http.StatusFound, response.Status)
	assert.Equal("/home", response.Headers.Get("Location"))
}

func TestRegisterBad(t *testing.T) {
	assert := assert.New(t)

	response := m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"name": "register", "email": "register@email.com", "confirm": "pass"},
		},
	)
	assert.Equal(http.StatusNotAcceptable, response.Status)
	assert.Contains(response.Body, "invalid username or password")

	response = m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"name": "register", "email": "register@email.com", "password": "pass", "confirm": "badpass"},
		},
	)
	assert.Equal(http.StatusNotAcceptable, response.Status)
	assert.Contains(response.Body, "passwords do not match")

	m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"name": "register", "email": "register@email.com", "password": "pass", "confirm": "pass"},
		},
	)
	response = m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"name": "register", "email": "register@email.com", "password": "pass", "confirm": "pass"},
		},
	)
	assert.Equal(http.StatusNotAcceptable, response.Status)
	assert.Contains(response.Body, "an account with this name or email address already exists")
}

func TestLoginPost(t *testing.T) {
	assert := assert.New(t)

	m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"name": "login", "email": "login@email.com", "password": "pass", "confirm": "pass"},
		},
	)

	response := m.POST(
		"/login",
		test.RequestOptions{
			Form: gin.H{"email": "login@email.com", "password": "pass"},
		},
	)

	assert.Equal(http.StatusFound, response.Status)
	assert.Equal("/home", response.Headers.Get("Location"))

	response = m.POST(
		"/login",
		test.RequestOptions{
			Form: gin.H{"email": "login@email.com", "password": "badpass"},
		},
	)

	assert.Equal(http.StatusUnauthorized, response.Status)
	assert.Contains(response.Body, "invalid username or password")
}

func TestLoadAuthorized(t *testing.T) {
	assert := assert.New(t)

	response := m.GET("/home")
	assert.Equal(http.StatusTemporaryRedirect, response.Status)
	assert.Equal("/login", response.Headers.Get("Location"))

	response = m.POST(
		"/register",
		test.RequestOptions{
			Form: gin.H{"name": "authorized", "email": "authorized@email.com", "password": "pass", "confirm": "pass"},
		},
	)

	response = m.WithSetCookie(response).GET("/home")

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "authorized")

	response = m.POST(
		"/login",
		test.RequestOptions{
			Form: gin.H{"email": "authorized@email.com", "password": "pass"},
		},
	)

	response = m.WithSetCookie(response).GET("/home")

	assert.Equal(http.StatusOK, response.Status)
	assert.Contains(response.Body, "authorized")
}
