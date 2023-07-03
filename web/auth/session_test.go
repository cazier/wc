package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cazier/wc/db/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// // SESSION

// // func TestAuthorized(t *testing.T) {
// // 	assert := assert.New(t)
// // }

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
	_, ok := GetUser(ctx)
	assert.False(ok)

	db.create("session", "session@email.com", "password")
	user := db.retrieve(models.User{Email: "session@email.com"})
	user.Session = models.NewToken("session", generateCookie())

	db.save(user)

	ctx.Request.AddCookie(&http.Cookie{Name: "session", Value: user.Session.Value})
	_, ok = GetUser(ctx)
	assert.True(ok)

	user.Session.CreatedAt = user.Session.CreatedAt.Add(-(2 * SESSION_LIFTIME))
	db.save(user)
	_, ok = GetUser(ctx)
	assert.False(ok)

	ctx.Request.AddCookie(&http.Cookie{Name: "session", Value: "invalid session token"})
	_, ok = GetUser(ctx)
	assert.False(ok)
}

func TestAddSessionCookie(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	assert.Empty(ctx.Writer.Header().Get("Set-Cookie"))

	cookie := generateCookie()

	addSessionCookie(ctx, cookie)
	assert.Contains(ctx.Writer.Header().Get("Set-Cookie"), cookie)
}

func TestEncodeParams(t *testing.T) {
	assert := assert.New(t)

	values := map[string]any{"a": "1", "b": "2"}
	assert.Equal("/encode?a=1&b=2", encodeParams("/encode", values))
	assert.Equal("/encode?a=1&b=2", encodeParams(encodeParams("/encode", values), values))

	assert.Equal("/encode/subdirectory?a+a=b+b", encodeParams("/encode/subdirectory", map[string]any{"a a": "b b"}))
}
