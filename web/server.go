package server

import (
	"github.com/cazier/wc/web/api"
	"github.com/cazier/wc/web/frontend"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func Init(database *gorm.DB) {
	gin.ForceConsoleColor()
	Engine = gin.Default()

	api.Init(database, Engine)
	frontend.Init(database, Engine)

	Engine.Run("0.0.0.0:1213")
}
