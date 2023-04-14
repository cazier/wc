package db

import (
	"log"
	"os"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/db/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			// SlowThreshold:             time.Second,   // Slow SQL threshold
			// LogLevel:                  logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			// // ParameterizedQueries:      true,          // Don't include params in the SQL log
			// Colorful: true, // Disable color
		},
	)

	// Globally mode
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Could not connect to the database")
	}

	database.AutoMigrate(&models.Country{}, &models.Match{})
	// load(database)
	LoadYaml(database)
	return database
}

func LoadYaml(database *gorm.DB) {
	teams, matches := utils.Import()

	storage := make(map[string]uint)

	for _, team := range teams {
		input := models.Country{Name: team.Name, FifaCode: team.Code, Group: team.Group}
		output := models.Country{}

		database.FirstOrCreate(&output, input)
		storage[team.Name] = output.ID
	}

	for _, match := range matches {
		input := models.Match{Day: 0, Played: false, AID: storage[match.A], BID: storage[match.B], Group: match.Group, When: match.Date}
		output := models.Match{}

		database.FirstOrCreate(&output, input)

		database.Model(storage[match.A]).Association("Matches").Append(&output)
		database.Model(storage[match.B]).Association("Matches").Append(&output)
	}

	addMatchDays(database)
}

func addMatchDays(database *gorm.DB) {
	countries := []models.Country{}
	database.Find(&countries)

	for _, country := range countries {
		matches := []models.Match{}

		database.Where(&models.Match{AID: country.ID}).Or(&models.Match{BID: country.ID}).Find(&matches).Order("when")

		for index, match := range matches {
			if match.Day == 0 {
				match.Day = uint(index + 1)
				database.Save(&match)
			}
		}

		// litter.Dump(matches)

	}

}
