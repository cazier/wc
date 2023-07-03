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
	g.POST("/login", loginPost, middlewares.Post())
	g.POST("/register", registerPost, middlewares.Post())

	r := g.Group("", middlewares.Authorized())
	r.GET("/", get)
	r.GET("/profile", get)
	r.POST("/profile", profilePost, middlewares.Post())
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
	c.Set(middlewares.RedirectKey, "/")
	auth.Create(c)
}

func loginPost(c *gin.Context) {
	c.Set(middlewares.RedirectKey, "/")
	auth.Login(c)
}

func profilePost(c *gin.Context) {
	c.Set(middlewares.RedirectKey, "/profile")
	auth.Update(c)
}

func postRoute(c *gin.Context, status int, target, fallback string, message map[string]any) {
	// switch status {
	// case http.StatusFound:
	// 	c.Redirect(status, target)
	// case http.StatusOK:
	// 	c.HTML(status, target, message)
	// case http.StatusInternalServerError:
	// 	c.HTML(status, "error.go.html", nil)
	// default:
	// 	c.HTML(status, fallback, message)
	// }
}
