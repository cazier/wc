package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB, engine *gin.Engine) {
	db = database
}
