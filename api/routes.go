package api

import (
	"github.com/cazier/wc/version"
	"github.com/gin-gonic/gin"
)

func getVersion(c *gin.Context) {
	c.JSON(200, gin.H{"version": version.Version})
}

func getPlayers(c *gin.Context) {
	c.JSON(500, gin.H{"error": "unkdnodwn"})
}
