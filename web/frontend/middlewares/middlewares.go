package middlewares

import (
	"fmt"
	"net/http"
	net "net/url"
	"path"

	"github.com/cazier/wc/web/exceptions"
	"github.com/gin-gonic/gin"
)

var TemplateKey = "Template"
var RedirectKey = "RedirectTarget"

func AssignTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(TemplateKey, TemplateName(c.FullPath()))
	}
}

func TemplateName(url string) string {
	parsed, _ := net.Parse(url)
	base := path.Base(parsed.Path)

	if base == "/" {
		base = "home"
	}

	return fmt.Sprintf("%s.go.html", base)
}

func Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			switch err {
			case exceptions.ErrAccountExists, exceptions.ErrAccountPasswordMismatch, exceptions.ErrAccountDetailsInvalid:
				c.HTML(http.StatusNotAcceptable, TemplateName("/"), gin.H{"message": err.Error()})
			case exceptions.ErrUnauthorized:
				c.HTML(http.StatusUnauthorized, TemplateName("/"), gin.H{"message": err.Error()})
			default:
				c.HTML(http.StatusInternalServerError, TemplateName("/error"), gin.H{})
			}
			c.Abort()
			return
		}

		c.Redirect(http.StatusFound, c.GetString(RedirectKey))
	}
}
