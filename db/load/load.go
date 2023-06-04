package load

import (
	"errors"
	"log"

	"github.com/cazier/wc/db"
	"github.com/cazier/wc/db/load/utils"
	"github.com/cazier/wc/db/models"
)

var cache map[string]models.Country

func Teams(path string) {
	var counter int64

	teams := utils.LoadTeams(path)

	teams = append([]utils.Team{
		{Name: "Team A", Group: "", Code: "<A>"},
		{Name: "Team B", Group: "", Code: "<B>"},
	}, teams...)

	for _, team := range teams {
		input := models.Country{Name: team.Name, FifaCode: team.Code, Group: team.Group}
		output := models.Country{}

		db.Database.FirstOrCreate(&output, input)

		counter++
	}
	log.Printf("Added %d (+2) countries to the database", counter-2)
}

func Matches(path string) {
	var counter int64
	cache, _ = cacheCountries()

	matches := utils.LoadMatches(path)

	for _, match := range matches {
		input := models.Match{
			AID:   cache[match.A].ID,
			BID:   cache[match.B].ID,
			When:  match.Date,
			Stage: match.Stage,
		}
		output := models.Match{}

		if db.Database.FirstOrCreate(&output, input).RowsAffected == 0 {
			continue
		}

		counter++
	}

	log.Printf("Added %d matches to the database", counter)

	db.AddMatchDays()
}

func Players(path string) {
	var counter int64
	cache, _ = cacheCountries()

	playerMap := utils.LoadPlayers(path)

	countryMap := make(map[string][]utils.Player)

	for _, player := range playerMap {
		countryMap[player.Country] = append(countryMap[player.Country], player)
	}

	for country, players := range countryMap {
		for _, player := range players {
			input := models.Player{
				Name:      player.Name,
				Position:  player.Position,
				Number:    player.Number,
				CountryID: cache[country].ID,
			}
			output := models.Player{}

			if db.Database.FirstOrCreate(&output, input).RowsAffected == 0 {
				continue
			}

			counter++
		}
	}
	log.Printf("Added %d players to the database", counter)
}

func cacheCountries() (map[string]models.Country, error) {
	var countries []models.Country

	if cache != nil {
		return cache, nil
	}

	tx := db.Database.Find(&countries)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, errors.New("cannot import match data when there are no countries in the table")
	}

	cache = make(map[string]models.Country)

	for _, item := range countries {
		cache[item.FifaCode] = item
		cache[item.Name] = item
	}

	log.Printf("Loaded %d countries into a cache map", len(cache))

	return cache, nil

}
