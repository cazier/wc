package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cazier/wc/db/models"

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

func AddMatchDays() {
	var teams []models.Country
	Database.Find(&teams)

	for _, team := range teams {
		matches := []models.Match{}

		Database.Where(&models.Match{AID: team.ID}).Or(&models.Match{BID: team.ID}).Find(&matches).Order("when")

		for index, match := range matches {
			if match.Day == 0 && match.Stage == models.GROUP {
				match.Day = uint(index + 1)
				match.Assigned = true
				Database.Save(&match)
			}
		}
	}
}
