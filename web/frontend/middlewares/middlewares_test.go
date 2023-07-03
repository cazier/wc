package middlewares

import (
	"os"
	"testing"

	"github.com/cazier/wc/db/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(tm *testing.M) {
	// test.NewMock(
	// 	&test.MockOptions{
	// 		Callback: func(database *gorm.DB, engine *gin.Engine) {
	// 			auth.Init(database, engine)
	// 			db = database
	// 		},
	// 		Models: []any{
	// 			&models.User{},
	// 		},
	// 	},
	// )

	os.Exit(tm.Run())
}

type MockedDb struct {
	mock.Mock
}

func (db *MockedDb) create(name, email, password string) (models.User, error) {
	args := db.Called(name, email, password)
	return models.User{Name: args.String(0)}, nil
}

func (db *MockedDb) retrieve(search models.User) models.User {
	args := db.Called(search)
	return models.User{Name: args.Get(0).(models.User).Name, Session: models.Token{Name: "session", Value: "session"}}
}

type mockedSession struct {
	mock.Mock
}

func (ms *mockedSession) GetUser(c *gin.Context) (models.User, bool) {
	ms.Called(c)
	return models.User{Name: "Alex"}, true
}

// func TestAuthorized(t *testing.T) {
// 	// assert := assert.New(t)

// 	auth := new(mockedSession)
// 	auth.On("GetUser", &gin.Context{}).Return(true)

// 	ctx := new(gin.Context)

// 	Authorized()(ctx)
// 	auth.AssertExpectations(t)

// 	// assert.Equal("Alex", user.Name)
// 	// assert.Nil(out)

// 	// mf := test.NewMockForm()
// 	// mf.Context.Request, _ = http.NewRequestWithContext(mf.Context, "GET", "", nil)
// 	// Authorized()(mf.Context)

// 	// assert.Zero(mf.Context.GetString(auth.StatusKey))
// 	// assert.Zero(mf.Context.GetString(auth.UserKey))

// 	// auth.Create(
// 	// 	mf.Form(
// 	// 		map[string]any{
// 	// 			"name":     "login",
// 	// 			"email":    "login@email.com",
// 	// 			"password": "password",
// 	// 			"confirm":  "password",
// 	// 		},
// 	// 	),
// 	// )

// 	// user := models.User{}
// 	// db.Where(&models.User{Name: "login"}).First(user)

// 	// mf.Context.Request.AddCookie(&http.Cookie{Name: "session", Value: user.Session.Value})

// 	// Authorized()(mf.Context)
// 	// assert.NotZero(mf.Context.GetString(auth.StatusKey))
// 	// assert.NotZero(mf.Context.GetString(auth.UserKey))
// }

func TestTemplateName(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("home.go.html", templateName("/"))
	assert.Equal("test.go.html", templateName("/test"))
	assert.Equal("test.go.html", templateName("/test?arbitrary=args"))
}
