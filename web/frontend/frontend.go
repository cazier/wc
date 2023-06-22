package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/cazier/wc/web/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var g *gin.Engine

//go:embed assets/*
var assets embed.FS

//go:embed templates/*
var templates embed.FS

func Init(database *gorm.DB, engine *gin.Engine) {
	_ = database
	g = engine

	addRoutes()
	loadStaticAssets()

	auth.Init(database, engine)
}

func addRoutes() {
	g.GET("/group", group)
	g.GET("/register", registerGet)
	g.POST("/register", registerPost)
	g.GET("/login", loginGet)
	g.POST("/login", loginPost)
	g.GET("/secret", secret)
}

func loadStaticAssets() {
	static, _ := fs.Sub(assets, "assets")
	g.StaticFS("/assets", http.FS(static))

	g.SetHTMLTemplate(template.Must(template.ParseFS(templates, "templates/*")))
}

func group(c *gin.Context) {
	c.HTML(200, "index.go.tmpl", map[string]any{})
}

func registerGet(c *gin.Context) {
	c.HTML(http.StatusOK, "register.go.tmpl", map[string]any{})
}

func registerPost(c *gin.Context) {
	status, message := auth.Create(c)

	switch status {
	case http.StatusFound:
		c.Redirect(status, "/secret")
	case http.StatusInternalServerError:
		c.HTML(status, "error.go.tmpl", nil)
	default:
		c.HTML(status, "register.go.tmpl", message)
	}
}

func secret(c *gin.Context) {
	c.HTML(http.StatusOK, "secret.go.tmpl", map[string]any{})
}

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.go.tmpl", map[string]any{})
}

func loginPost(c *gin.Context) {
	status, message := auth.Login(c)

	switch status {
	case http.StatusFound:
		c.Redirect(status, "/secret")
	case http.StatusInternalServerError:
		c.HTML(status, "error.go.tmpl", nil)
	default:
		c.HTML(status, "login.go.tmpl", message)
	}
}
