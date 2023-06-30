package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	test "github.com/cazier/wc/testing"
	"github.com/cazier/wc/web/exceptions"
	"github.com/gin-gonic/gin"

	"github.com/cazier/wc/db/models"
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

func TestDbCreate(t *testing.T) {
	assert := assert.New(t)

	_, err := create("dbcreate", "dbcreate@email.com", "password")
	assert.Nil(err)

	_, err = create("new_dbcreate", "dbcreate@email.com", "passwordpassword")
	assert.Errorf(err, "an account with this name or email address already exists")

	_, err = create("dbcreate", "new_dbcreate@email.com", "passwordpassword")
	assert.Errorf(err, "an account with this name or email address already exists")

	_, err = create("dbcreate", "dbcreate@email.com", "password")
	assert.Errorf(err, "an account with this name or email address already exists")
}

func TestRetrieve(t *testing.T) {
	assert := assert.New(t)

	user := retrieve(models.User{Email: "retrieve@email.com"})
	assert.Nil(user.Salt)
	assert.Nil(user.Hash)

	user, _ = create("retrieve", "retrieve@email.com", "password")
	assert.NotNil(user.Salt)
	assert.NotNil(user.Hash)

	assert.False(user.IsNil())

	retrieved := retrieve(models.User{Email: "retrieve@email.com"})
	assert.EqualExportedValues(user, retrieved)
}

func TestIsValid(t *testing.T) {
	assert := assert.New(t)

	create("valid", "isvalid@email.com", "password")

	assert.True(isValid("isvalid@email.com", "password"))

	assert.False(isValid("isvalid@email.com", "wrongpassword"))
	assert.False(isValid("isinvalid@email.com", "password"))
}

func TestGenerateSalt(t *testing.T) {
	assert := assert.New(t)

	salt := generateSalt()

	assert.NotNil(salt)
	assert.Len(salt, saltLength)
}

func TestGenerateHash(t *testing.T) {
	assert := assert.New(t)

	hash := generateHash("password", generateSalt())

	assert.NotNil(hash)
	assert.Len(hash, hashLength)
}

func TestGenerateCookie(t *testing.T) {
	assert := assert.New(t)

	session := generateCookie()
	assert.NotZero(session)
	assert.GreaterOrEqual(len(session), sessionCookieLength)

	assert.NotEqual(session, generateCookie())
}

func TestGenerateCsrf(t *testing.T) {
	assert := assert.New(t)

	session := generateCsrf()
	assert.NotZero(session)
	assert.GreaterOrEqual(len(session), csrfLength)

	assert.NotEqual(session, generateCsrf())
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{Header: make(http.Header)}
	_, ok := getUser(ctx)
	assert.False(ok)

	create("session", "session@email.com", "password")
	user := retrieve(models.User{Email: "session@email.com"})
	user.Session = models.NewToken("session", generateCookie())

	save(user)

	ctx.Request.AddCookie(&http.Cookie{Name: "session", Value: user.Session.Value})
	_, ok = getUser(ctx)
	assert.True(ok)

	user.Session.CreatedAt = user.Session.CreatedAt.Add(-(2 * SESSION_LIFTIME))
	save(user)
	_, ok = getUser(ctx)
	assert.False(ok)

	ctx.Request.AddCookie(&http.Cookie{Name: "session", Value: "invalid session token"})
	_, ok = getUser(ctx)
	assert.False(ok)
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	mf := test.MockForm{}
	type msa map[string]any

	Create(mf.Form(msa{"name": "create", "email": "create@email.com", "password": "pass", "confirm": "pass"}))
	assert.Empty(mf.Errors())

	Create(mf.Form(msa{"email": "create@email.com", "confirm": "password"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountDetailsInvalid.Error())

	Create(mf.Form(msa{"name": "create", "email": "create@email.com", "password": "pass", "confirm": "badpass"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountPasswordMismatch.Error())

	Create(mf.Form(msa{"name": "create", "email": "create@email.com", "password": "pass", "confirm": "pass"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountExists.Error())
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	mf := test.MockForm{}
	type msa map[string]any

	Create(mf.Form(msa{"email": "login@email.com", "password": "password", "confirm": "password"}))

	Login(mf.Form(msa{"email": "login@email.com", "password": "password"}))
	assert.Empty(mf.Errors())

	Login(mf.Form(msa{"email": "login@email.com"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountDetailsInvalid.Error())

	Login(mf.Form(msa{"email": "login@email.com", "password": "wrongpassword"}))
	assert.Contains(mf.Errors(), exceptions.ErrUnauthorized.Error())
}

func TestEncodeQueries(t *testing.T) {
	assert := assert.New(t)

	values := map[string]any{"a": "1", "b": "2"}
	assert.Equal("/encode?a=1&b=2", encodeParams("/encode", values))
	assert.Equal("/encode?a=1&b=2", encodeParams(encodeParams("/encode", values), values))

	assert.Equal("/encode/subdirectory?a+a=b+b", encodeParams("/encode/subdirectory", map[string]any{"a a": "b b"}))
}
