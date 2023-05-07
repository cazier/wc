package api

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(g *gin.Engine) {
	g.GET("/version", getVersion)
	g.GET("/players", getPlayers)
}

func Init() {
	Api := gin.Default()
	setupRoutes(Api)
	Api.Run("0.0.0.0:1213")
}
