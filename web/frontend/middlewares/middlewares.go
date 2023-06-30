package middlewares

import (
	"fmt"
	net "net/url"
	"path"

	"github.com/gin-gonic/gin"
)

var TemplateKey = "Template"

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
	}
}
