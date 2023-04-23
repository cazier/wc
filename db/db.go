package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/db/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func Init(databasePath string) {
	var err error

	var logPath = filepath.Dir(databasePath)

	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

	newLogger := logger.New(
		log.New(logFile, "db.db", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Do include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	defer logFile.Close()

	Database, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatalf("Could not connect to the database. %s", err)
	}
}

// Create tables in the database for the specific data models
func LinkTables() {
	Database.AutoMigrate(
		&models.Country{},
		&models.Player{},
		&models.Match{},
		&models.MatchResult{},
	)
}

// Import yaml file data into the database
func LoadYaml(importTeams, importMatches string) {
	var counter int64
	storage := make(map[string]uint)

	if importTeams != "" {
		teams := utils.LoadTeams(importTeams)

		teams = append([]utils.Team{
			{Name: "Team A", Group: "", Code: "<A>"},
			{Name: "Team B", Group: "", Code: "<B>"},
		}, teams...)

		for _, team := range teams {
			input := models.Country{Name: team.Name, FifaCode: team.Code, Group: team.Group}
			output := models.Country{}

			Database.FirstOrCreate(&output, input)
			storage[team.Name] = output.ID

			counter++
		}
		log.Printf("Added %d countries to the database", counter)

	} else {
		var teams []models.Country
		tx := Database.Find(&teams)

		if tx.Error != nil {
			log.Fatalf("An error occurred with the database. %s", tx.Error)
		}

		if tx.RowsAffected == 0 {
			log.Fatal("Cannot import match data when there are no countries in the table.")
		}

		for _, team := range teams {
			storage[team.Name] = team.ID
		}

		log.Printf("Retrieved %d countries from the database", len(teams))
	}

	counter = 0

	if importMatches != "" {
		matches := utils.LoadMatches(importMatches)

		for _, match := range matches {
			input := models.Match{AID: storage[match.A], BID: storage[match.B], When: match.Date, Stage: match.Stage}
			output := models.Match{}

			if Database.FirstOrCreate(&output, input).RowsAffected == 0 {
				break
			}

			counter++

			Database.Model(storage[match.A]).Association("Matches").Append(&output)
			Database.Model(storage[match.B]).Association("Matches").Append(&output)
		}
		log.Printf("Added %d matches to the database", counter)

		addMatchDays()
	}

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
