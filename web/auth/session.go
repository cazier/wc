package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cazier/wc/db/models"
	"github.com/gin-gonic/gin"
)

// var secret = []byte("correcthorsebatterystaple")

const UserKey = "User"
const AuthStatusKey = "Authorization"

const csrfLength = 64
const sessionCookieLength = 128

const SESSION_LIFTIME = 3 * 24 * time.Hour

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, ok := getUser(c); !ok {
			//TODO: Store the desired page internally and load that, as desired?
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			// target := encodeParams("/login", map[string]any{"redirect": c.FullPath()})
			// c.Redirect(http.StatusFound, target)
			c.Abort()
		} else {
			c.Set(AuthStatusKey, http.StatusAccepted)
			c.Set(UserKey, user)
		}
	}
}

func generateCookie() string {
	id := make([]byte, sessionCookieLength)
	rand.Read(id)

	return strings.ReplaceAll(base64.URLEncoding.EncodeToString(id), "=", "")
}

func generateCsrf() string {
	id := make([]byte, csrfLength)
	rand.Read(id)

	return base64.URLEncoding.EncodeToString(id)
}

func getUser(c *gin.Context) (models.User, bool) {
	session, err := c.Cookie("session")
	blank := models.User{}

	if err != nil {
		return blank, false
	}

	user := retrieve(models.User{Session: models.Token{Value: session}})

	if user.IsNil() || !user.Session.IsValid(SESSION_LIFTIME) {
		return blank, false
	}

	return user, true
}

func addSessionCookie(c *gin.Context, cookie string) {
	c.SetCookie("session", cookie, int(SESSION_LIFTIME.Seconds()), "", "", false, true)
}

func encodeParams(path string, params map[string]any) string {
	encode, _ := url.Parse(path)
	v, _ := url.ParseQuery(encode.RawQuery)

	for key, value := range params {
		if v.Has(key) {
			continue
		}
		v.Add(key, value.(string))
	}

	encode.RawQuery = v.Encode()

	return encode.String()
}
