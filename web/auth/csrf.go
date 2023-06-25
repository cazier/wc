package auth

import (
	"net/http"
	"time"

	"github.com/cazier/wc/db/models"
	"github.com/gin-gonic/gin"
)

var storage map[string]models.Token

func init() {
	storage = make(map[string]models.Token)
}

const CsrfKey = "csrf"
const csrf_lifetime = time.Minute * 10

func IncludeCsrfToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := generateCsrf()

		if session, err := c.Cookie("session"); err == nil {
			user := retrieve(models.User{Session: models.Token{Value: session}})

			if user.IsNil() {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			user.Csrf = models.NewToken("csrf", token)
			c.Set(CsrfKey, token)
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func GetToken(c *gin.Context) string {
	return storage["infantino@fifa.com"].Value
}

func validCsrf(value string) bool {
	if csrf, ok := storage[value]; ok {
		if csrf.IsValid(csrf_lifetime) {
			return true
		}
	}
	return false
}
