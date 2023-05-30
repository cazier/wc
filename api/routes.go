package api

import (
	"database/sql"

	"github.com/cazier/wc/api/exceptions"
	"github.com/cazier/wc/db"
	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/version"
	"github.com/gin-gonic/gin"
)

func getVersion(c *gin.Context) {
	c.JSON(200, gin.H{"version": version.Version})
}

func getPlayer(c *gin.Context) {
	if resp, ok := queryPlayers(c, false); ok {
		c.JSON(200, gin.H{"data": resp[0]})
	}
}

func getPlayers(c *gin.Context) {
	if resp, ok := queryPlayers(c, true); ok {
		c.JSON(200, gin.H{"data": resp})
	}
}

func getCountry(c *gin.Context) {
	if resp, ok := queryCountries(c, false); ok {
		c.JSON(200, gin.H{"data": resp[0]})
	}
}

func getCountries(c *gin.Context) {
	if resp, ok := queryCountries(c, true); ok {
		c.JSON(200, gin.H{"data": resp})
	}
}

func getMatch(c *gin.Context) {
	if resp, ok := queryMatches(c, false); ok {
		c.JSON(200, gin.H{"data": resp[0]})
	}
}
func getMatches(c *gin.Context) {
	if resp, ok := queryMatches(c, true); ok {
		c.JSON(200, gin.H{"data": resp})
	}
}

func getPlayerMatches(c *gin.Context) {
	var matches []models.Match
	var search models.Player

	err := c.ShouldBindUri(&search)
	if exceptions.JsonResponse(c, err) {
		return
	}

	for _, param := range c.Params {
		if param.Value == "" {
			exceptions.JsonResponse(c, &exceptions.InvalidValueError{})
			return
		}
	}

	db.Database.Joins("ACountry").Joins("BCountry").
		Joins("JOIN players ON `players`.`country_id` = a_id OR `players`.`country_id` = b_id").
		Where(
			"`players`.`name` LIKE @name OR `players`.`id` = @id",
			sql.Named("name", search.Name),
			sql.Named("id", search.ID),
		).
		Order("`matches`.`when`").
		Find(&matches)

	if len(matches) == 0 {
		exceptions.JsonResponse(c, &exceptions.NoResultsFoundError{})
		return
	}

	c.JSON(200, gin.H{"data": matches})
}

func getCountryMatches(c *gin.Context) {
	var matches []models.Match
	var search models.Country

	err := c.ShouldBindUri(&search)
	if exceptions.JsonResponse(c, err) {
		return
	}

	for _, param := range c.Params {
		if param.Value == "" {
			exceptions.JsonResponse(c, &exceptions.InvalidValueError{})
			return
		}
	}

	// TODO clean this up
	db.Database.Joins("ACountry").Joins("BCountry").
		Where("`ACountry`.`Name` LIKE @name OR `BCountry`.`Name` LIKE @name", sql.Named("name", search.Name)).
		Or("`ACountry`.`ID` = @id OR `BCountry`.`ID` = @id", sql.Named("id", search.ID)).
		Order("`matches`.`when`").
		Find(&matches)

	if len(matches) == 0 {
		exceptions.JsonResponse(c, &exceptions.NoResultsFoundError{})
		return
	}

	c.JSON(200, gin.H{"data": matches})
}

func getCountryPlayers(c *gin.Context) {
	type Player struct {
		ID uint `json:"id"`

		Name     string `json:"name"`
		Position string `json:"position"`
		Number   int    `json:"number"`

		Goals  uint `json:"goals" `
		Yellow uint `json:"yellows"`
		Red    uint `json:"reds"`
		Saves  int  `json:"saves"`
	}

	var players []Player
	var search models.Country

	err := c.ShouldBindUri(&search)
	if exceptions.JsonResponse(c, err) {
		return
	}

	for _, param := range c.Params {
		if param.Value == "" {
			exceptions.JsonResponse(c, &exceptions.InvalidValueError{})
			return
		}
	}

	// TODO clean this up
	db.Database.Model(&models.Player{}).Joins("Country").
		Where(
			"Country.Name = @name OR Country.ID = @id",
			sql.Named("name", search.Name),
			sql.Named("id", search.ID),
		).Find(&players)

	if len(players) == 0 {
		exceptions.JsonResponse(c, &exceptions.NoResultsFoundError{})
		return
	}

	c.JSON(200, gin.H{"data": players})
}
