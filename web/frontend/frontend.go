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
	auth.Create(c)

	// switch status {
	// case http.StatusAccepted:
	// 	c.HTML(status, "secret.go.tmpl", map[string]any{})
	// case http.StatusNotAcceptable, http.StatusUnauthorized:
	// 	c.HTML(status, "login.go.tmpl", map[string]any{})
	// default:
	// 	c.HTML(http.StatusInternalServerError, "error.go.tmpl", map[string]any{})
	// }
}

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.go.tmpl", map[string]any{})
}

func loginPost(c *gin.Context) {
	status := auth.Login(c)

	switch status {
	case http.StatusAccepted:
		c.HTML(status, "secret.go.tmpl", map[string]any{})
	case http.StatusNotAcceptable, http.StatusUnauthorized:
		c.HTML(status, "login.go.tmpl", map[string]any{})
	default:
		c.HTML(http.StatusInternalServerError, "error.go.tmpl", map[string]any{})
	}
}
