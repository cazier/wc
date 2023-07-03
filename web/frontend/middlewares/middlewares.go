package middlewares

import (
	"fmt"
	"net/http"
	net "net/url"
	"path"

	"github.com/cazier/wc/web/auth"
	"github.com/cazier/wc/web/exceptions"
	"github.com/gin-gonic/gin"
)

var TemplateKey = "Template"
var RedirectKey = "RedirectTarget"

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, ok := auth.GetUser(c); !ok {
			//TODO: Store the desired page internally and load that, as desired?
			// target := encodeParams("/login", map[string]any{"redirect": c.FullPath()})
			// c.Redirect(http.StatusFound, target)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		} else {
			c.Set(auth.StatusKey, http.StatusAccepted)
			c.Set(auth.UserKey, user)
		}
	}
}

func AssignTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(TemplateKey, templateName(c.FullPath()))
	}
}

func Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			switch err {
			case exceptions.ErrAccountExists, exceptions.ErrAccountPasswordMismatch, exceptions.ErrAccountDetailsInvalid:
				c.HTML(http.StatusNotAcceptable, templateName("/"), gin.H{"message": err.Error()})
			case exceptions.ErrUnauthorized:
				c.HTML(http.StatusUnauthorized, templateName("/"), gin.H{"message": err.Error()})
			default:
				c.HTML(http.StatusInternalServerError, templateName("/error"), gin.H{})
			}
			c.Abort()
			return
		}

		c.Redirect(http.StatusFound, c.GetString(RedirectKey))
	}
}

func templateName(url string) string {
	parsed, _ := net.Parse(url)
	base := path.Base(parsed.Path)

	if base == "/" {
		base = "home"
	}

	return fmt.Sprintf("%s.go.html", base)
}
