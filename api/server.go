package api

import (
	"github.com/gin-gonic/gin"
)

var Api *gin.Engine

func Init() {
	gin.ForceConsoleColor()

	Api = gin.Default()
	setupRoutes(Api)
	Api.Run("0.0.0.0:1213")
}

func setupRoutes(g *gin.Engine) {
	utilities(g)
	matches(g)
	players(g)
	countries(g)
}

func utilities(g *gin.Engine) {
	g.GET("/version", getVersion)
}

func players(g *gin.Engine) {
	g.GET("/player", getPlayers)

	g.GET("/player/id/:id", getPlayer)
	g.GET("/player/name/:name", getPlayer)
	// TODO: Positions? Numbers? Goals? Cards?

	g.GET("/country/id/:id/players", getCountryPlayers)
	g.GET("/country/name/:name/players", getCountryPlayers)
}

func countries(g *gin.Engine) {
	g.GET("/country", getCountries)

	g.GET("/country/id/:id", getCountry)
	g.GET("/country/code/:code", getCountry)
	g.GET("/country/name/:name", getCountry)

	g.GET("/country/group/:group", getCountries)
}

func matches(g *gin.Engine) {
	g.GET("/player/id/:id/matches", getPlayerMatches)
	g.GET("/player/name/:name/matches", getPlayerMatches)

	g.GET("/country/id/:id/matches", getCountryMatches)
	g.GET("/country/name/:name/matches", getCountryMatches)

	g.GET("/match", getMatches)
	g.GET("/match/id/:id", getMatch)
	// g.GET("/match/between/:country_a/:country_b", getMatch)

	g.GET("/match/day/:day", getMatches)
	g.GET("/match/group/:group", getMatches)
	g.GET("/match/stage/:stage", getMatches)
}
