package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/cazier/wc/db/models"
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
	g.GET("/login", loginGet)
	g.GET("/register", registerGet)
	g.POST("/login", loginPost)
	g.POST("/register", registerPost)

	r := g.Group("", auth.Authorized())
	r.GET("/home", home)
	r.GET("/group", group)
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
		c.Redirect(status, "/home")
	case http.StatusInternalServerError:
		c.HTML(status, "error.go.tmpl", nil)
	default:
		c.HTML(status, "register.go.tmpl", message)
	}
}

func home(c *gin.Context) {
	status, _ := c.Get(auth.AuthStatusKey)
	user, _ := c.Get(auth.UserKey)

	switch status.(int) {
	case http.StatusAccepted:
		c.HTML(http.StatusOK, "home.go.tmpl", user.(models.User).Serialize())
	case http.StatusUnauthorized:
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	default:
		c.HTML(http.StatusInternalServerError, "error.go.tmpl", nil)
	}
}

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.go.tmpl", gin.H{})
}

func loginPost(c *gin.Context) {
	status, message := auth.Login(c)

	switch status {
	case http.StatusFound:
		c.Redirect(status, "/home")
	case http.StatusInternalServerError:
		c.HTML(status, "error.go.tmpl", nil)
	default:
		c.HTML(status, "login.go.tmpl", message)
	}
}
