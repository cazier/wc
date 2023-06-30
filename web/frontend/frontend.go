package frontend

import (
	"embed"
	"fmt"
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
	g.Use()

	g.GET("/login", loginGet)
	g.GET("/register", registerGet)
	g.POST("/login", loginPost)
	g.POST("/register", registerPost)

	r := g.Group("", auth.Authorized())
	r.GET("/", home)
	r.GET("/group", group)
	r.GET("/profile", profileGet)
	r.POST("/profile", profilePost)
}

func loadStaticAssets() {
	static, _ := fs.Sub(assets, "assets")
	g.StaticFS("/assets", http.FS(static))

	g.SetHTMLTemplate(template.Must(template.ParseFS(templates, "templates/*")))
}

func group(c *gin.Context) {
	c.HTML(200, "index.go.html", map[string]any{})
}

func registerGet(c *gin.Context) {
	c.HTML(http.StatusOK, "register.go.html", gin.H{"user": nil})
}

func registerPost(c *gin.Context) {
	status, message := auth.Create(c)
	postRoute(c, status, "/", "register.go.html", message)
}

func home(c *gin.Context) {
	status, _ := c.Get(auth.AuthStatusKey)
	user, _ := c.Get(auth.UserKey)
	postRoute(c, status.(int), "home.go.html", "/", gin.H{"user": user.(models.User).Serialize()})
}

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.go.html", gin.H{"user": nil})
}

func loginPost(c *gin.Context) {
	status, message := auth.Login(c)
	postRoute(c, status, "/", "login.go.html", message)
}

func profileGet(c *gin.Context) {
	user, _ := c.Get(auth.UserKey)
	fmt.Println(user)
	c.HTML(http.StatusOK, "profile.go.html", gin.H{"user": nil})
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
