package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db AuthenticationDatabase
var Auth Authentication

func Init(database *gorm.DB, engine *gin.Engine) {
	db = &RuntimeAuthenticationDatabase{db: database}
	Auth = new(RuntimeAuthentication)
}

func Create(c *gin.Context) {
	Auth.Create(c)
}

func Login(c *gin.Context) {
	Auth.Login(c)
}

func Update(c *gin.Context) {
	Auth.Update(c)
}
