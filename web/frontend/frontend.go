package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/web/auth"
	"github.com/cazier/wc/web/frontend/middlewares"
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
	g.Use(middlewares.AssignTemplate())

	g.GET("/login", get)
	g.GET("/register", get)
	g.POST("/login", loginPost)
	g.POST("/register", registerPost)

	r := g.Group("", auth.Authorized())
	r.GET("/", get)
	r.GET("/profile", get)
	r.POST("/profile", profilePost)
}

func loadStaticAssets() {
	static, _ := fs.Sub(assets, "assets")
	g.StaticFS("/assets", http.FS(static))

	g.SetHTMLTemplate(template.Must(template.ParseFS(templates, "templates/*")))
}

func get(c *gin.Context) {
	var data map[string]string
	user, exists := c.Get(auth.UserKey)

	if exists {
		data = user.(models.User).Serialize()
	}

	c.HTML(http.StatusOK, c.GetString(middlewares.TemplateKey), gin.H{"user": data})
}

func registerPost(c *gin.Context) {
	status, message := auth.Create(c)
	postRoute(c, status, "/", "register.go.html", message)
}

func loginPost(c *gin.Context) {
	status, message := auth.Login(c)
	postRoute(c, status, "/", "login.go.html", message)
}

func profilePost(c *gin.Context) {
	status, message := auth.Update(c)
	postRoute(c, status, "/", "profile.go.html", message)
}

func postRoute(c *gin.Context, status int, target, fallback string, message map[string]any) {
	switch status {
	case http.StatusFound:
		c.Redirect(status, target)
	case http.StatusOK:
		c.HTML(status, target, message)
	case http.StatusInternalServerError:
		c.HTML(status, "error.go.html", nil)
	default:
		c.HTML(status, fallback, message)
	}
}
