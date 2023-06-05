package api

import (
	"github.com/gin-gonic/gin"
)

var group *gin.RouterGroup

func Init(g *gin.Engine) {
	group = g.Group("/api")

	players()
	countries()
	matches()
	utilities()
}

func players() {
	group.GET("/player", getPlayers)

	group.GET("/player/id/:id", getPlayer)
	group.GET("/player/name/:name", getPlayer)
	// TODO: Positions? Numbers? Goals? Cards?

	group.GET("/country/id/:id/players", getCountryPlayers)
	group.GET("/country/name/:name/players", getCountryPlayers)
}

func countries() {
	group.GET("/country", getCountries)

	group.GET("/country/id/:id", getCountry)
	group.GET("/country/code/:code", getCountry)
	group.GET("/country/name/:name", getCountry)

	group.GET("/country/group/:group", getCountries)
}

func matches() {
	group.GET("/player/id/:id/matches", getPlayerMatches)
	group.GET("/player/name/:name/matches", getPlayerMatches)

	group.GET("/country/id/:id/matches", getCountryMatches)
	group.GET("/country/name/:name/matches", getCountryMatches)

	group.GET("/match", getMatches)
	group.GET("/match/id/:id", getMatch)
	// group.GET("/match/between/:country_a/:country_b", getMatch)

	group.GET("/match/day/:day", getMatches)
	group.GET("/match/group/:group", getMatches)
	group.GET("/match/stage/:stage", getMatches)
}

func utilities() {
	group.GET("/version", getVersion)
}
