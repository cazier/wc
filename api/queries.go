package api

import (
	"database/sql"

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
	if options.callback != nil {
		options.query = options.callback(options.query)
	}

	bindings, err := bindUri(c, &search)
	if exceptions.JsonResponse(c, err) {
		return nil, false
	}

	if bindings {
		options.query = options.query.Where(&search)
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

func adaptNameCase(c *gin.Context) (bool, string) {
	for index, param := range c.Params {
		if param.Key == "name" && param.Value != "" {
			c.Params[index].Value = ""
			return true, param.Value
		}
	}
	return false, ""
}

func queryPlayers(c *gin.Context, multiple bool) ([]models.Player, bool) {
	// TODO special characters in player name?
	var players []models.Player
	var search models.Player

	options := QueryOptions{query: db.Database.Joins("Country"), multiple: multiple}

	if swap, name := adaptNameCase(c); swap {
		options.callback = func(tx *gorm.DB) *gorm.DB {
			return tx.Where("`players`.`name` LIKE @name", sql.Named("name", name))
		}
	}

	return query(search, players, c, &options)
}

func queryCountries(c *gin.Context, multiple bool) ([]models.Country, bool) {
	var countries []models.Country
	var search models.Country

	options := QueryOptions{query: db.Database, multiple: multiple}

	if multiple {
		// Ignore the `Team A` and `Team B` placeholder teams
		options.query = db.Database.Where("`countries`.`fifa_code` <> \"<A>\" AND `countries`.`fifa_code` <> \"<B>\"")
	}

	if swap, name := adaptNameCase(c); swap {
		options.callback = func(tx *gorm.DB) *gorm.DB {
			return tx.Where("`countries`.`name` LIKE @name", sql.Named("name", name))
		}
	}

	return query(search, countries, c, &options)
}

func queryMatches(c *gin.Context, multiple bool) ([]models.Match, bool) {
	var matches []models.Match
	var search models.Match

	tx := db.Database.Joins("ACountry").Joins("BCountry").Order("`matches`.`when`")

	if group, found := c.Params.Get("group"); found {
		tx = tx.Where("`ACountry`.`group` LIKE @group OR `BCountry`.`group` LIKE @group", sql.Named("group", group))
	}

	return query(search, matches, c, &QueryOptions{query: tx, multiple: multiple})
}
