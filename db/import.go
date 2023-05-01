package db

import (
	"log"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/db/utils"
)

var cache map[string]uint

func LoadTeams(path string) {
	var counter int64

	teams := utils.LoadTeams(path)

	teams = append([]utils.Team{
		{Name: "Team A", Group: "", Code: "<A>"},
		{Name: "Team B", Group: "", Code: "<B>"},
	}, teams...)

	for _, team := range teams {
		input := models.Country{Name: team.Name, FifaCode: team.Code, Group: team.Group}
		output := models.Country{}

		Database.FirstOrCreate(&output, input)

		counter++
	}
	log.Printf("Added %d (+2) countries to the database", counter-2)
}

func LoadMatches(path string) {
	var counter int64
	cache = cacheCountries()

	matches := utils.LoadMatches(path)

	for _, match := range matches {
		input := models.Match{AID: cache[match.A], BID: cache[match.B], When: match.Date, Stage: match.Stage}
		output := models.Match{}

		if Database.FirstOrCreate(&output, input).RowsAffected == 0 {
			continue
		}

		Database.First(&models.Country{}, cache[match.A]).Association("Matches").Append(&output)
		Database.First(&models.Country{}, cache[match.B]).Association("Matches").Append(&output)

		counter++
	}

	log.Printf("Added %d matches to the database", counter)

	addMatchDays()
}

func LoadPlayers(path string) {
	var counter int64
	cache = cacheCountries()

	playerMap := utils.LoadPlayers(path)

	for country, players := range playerMap {
		for _, player := range players {
			input := models.Player{Name: player.Name, Position: player.Postition, Number: player.Number, CountryID: cache[country]}
			output := models.Player{}

			if Database.FirstOrCreate(&output, input).RowsAffected == 0 {
				continue
			}

			Database.First(&models.Country{}, cache[country]).Association("Players").Append(&output)

			counter++
		}
	}
	log.Printf("Added %d players to the database", counter)
}

func cacheCountries() map[string]uint {
	var countries []models.Country

	if cache != nil {
		return cache
	}

	cache = make(map[string]uint)

	tx := Database.Find(&countries)

	if tx.Error != nil {
		log.Fatalf("An error occurred with the database. %s", tx.Error)
	}

	if tx.RowsAffected == 0 {
		log.Fatal("Cannot import match data when there are no countries in the table.")
	}

	for _, item := range countries {
		cache[item.Name] = item.ID
	}

	log.Printf("Loaded %d countries into a cache map", len(cache))

	return cache

}

func addMatchDays() {
	var teams []models.Country
	Database.Find(&teams)

	for _, team := range teams {
		matches := []models.Match{}

		Database.Where(&models.Match{AID: team.ID}).Or(&models.Match{BID: team.ID}).Find(&matches).Order("when")

		for index, match := range matches {
			if match.Day == 0 && match.Stage == utils.GROUP {
				match.Day = uint(index + 1)
				match.Assigned = true
				Database.Save(&match)
			}
		}
	}
}
