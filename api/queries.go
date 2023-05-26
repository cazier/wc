package api

import (
	"github.com/cazier/wc/api/exceptions"
	"github.com/cazier/wc/db"
	"github.com/cazier/wc/db/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QueryOptions struct {
	query    *gorm.DB
	multiple bool
	callback func(tx *gorm.DB) *gorm.DB
}

func query[M any](search M, dest []M, c *gin.Context, options *QueryOptions) ([]M, bool) {
	bindings, err := bindUri(c, &search)
	if exceptions.JsonResponse(c, err) {
		return nil, false
	}

	if bindings {
		options.query = options.query.Where(&search)
	}

	if options.callback != nil {
		options.query = options.callback(options.query)
	}

	if options.multiple {
		options.query.Find(&dest)
	} else {
		options.query.First(&dest)
	}

	if len(dest) == 0 {
		exceptions.JsonResponse(c, &exceptions.NoResultsFoundError{})
		return nil, false
	}

	return dest, true
}

func bindUri(c *gin.Context, obj any) (bool, error) {
	if c.Params == nil {
		return false, nil
	}

	if id, found := c.Params.Get("id"); found && id == "0" {
		return false, &exceptions.NoResultsFoundError{}
	}

	return true, c.ShouldBindUri(obj)
}

func queryPlayers(c *gin.Context, multiple bool) ([]models.Player, bool) {
	// TODO special characters in player name?
	var players []models.Player
	var search models.Player

	return query(search, players, c, &QueryOptions{query: db.Database.Joins("Country"), multiple: multiple})
}

func queryCountries(c *gin.Context, multiple bool) ([]models.Country, bool) {
	var countries []models.Country
	var search models.Country

	tx := db.Database

	if multiple {
		tx = tx.Where("`countries`.`group` <> ?", "")
	}

	return query(search, countries, c, &QueryOptions{query: tx, multiple: multiple})
}

func queryMatches(c *gin.Context, multiple bool) ([]models.Match, bool) {
	var matches []models.Match
	var search models.Match

	tx := db.Database.Joins("ACountry").Joins("BCountry")

	if group, found := c.Params.Get("group"); found {
		tx = tx.Where("`ACountry`.`group` = ? OR `BCountry`.`group` = ?", group, group)
	}

	return query(search, matches, c, &QueryOptions{query: tx, multiple: multiple})
}
