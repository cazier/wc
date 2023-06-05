package server

import (
	"github.com/cazier/wc/web/api"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func Init() {
	gin.ForceConsoleColor()

	Engine = gin.Default()
	Engine.Run("0.0.0.0:1213")

	api.Init(Engine)
}
