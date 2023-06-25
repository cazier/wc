package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cazier/wc/db/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func open(database gorm.Dialector, logLevel int, logPath string, purge Purge) {
	var err error
	var writer *log.Logger

	if logPath == "stdout" {
		writer = log.New(os.Stdout, "\n", log.LstdFlags)
	} else {
		logPath = filepath.Dir(logPath)

		logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		writer = log.New(logFile, "db.db", log.LstdFlags)

		defer logFile.Close()
	}

	newLogger := logger.New(
		writer,
		logger.Config{
			SlowThreshold:             time.Second,               // Slow SQL threshold
			LogLevel:                  logger.LogLevel(logLevel), // Log level
			IgnoreRecordNotFoundError: true,                      // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,                     // Do include params in the SQL log
			Colorful:                  false,                     // Disable color
		},
	)

	Database, err = gorm.Open(database, &gorm.Config{Logger: newLogger, TranslateError: true})

	if err != nil {
		log.Fatalf("Could not connect to the database. %s", err)
	}

	if purge != 0 {
		LinkTables(purge == 1)
	}
}

func InitMariaDB(options *MariaDBOptions) {
	options.validate()

	dialect := mysql.Open(
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?%s",
			options.Username,
			options.Password,
			options.Host,
			options.Port,
			options.Database,
			options.Other,
		),
	)

	open(dialect, options.LogLevel, options.LogPath, options.Purge)
}

func InitSqlite(options *SqliteDBOptions) {
	options.validate()

	dialect := sqlite.Open(fmt.Sprintf("%s?%s", options.Path, options.Other))
	open(dialect, options.LogLevel, options.LogPath, options.Purge)
}

// Create tables in the database for the specific data models
func LinkTables(purge bool) {
	if purge {
		Database.Migrator().DropTable(
			&models.Country{},
			&models.Player{},
			&models.Match{},
			&models.User{},
			// &models.MatchResult{},
		)
	}
	Database.AutoMigrate(
		&models.Country{},
		&models.Player{},
		&models.Match{},
		&models.User{},
		// &models.MatchResult{},
	)
}

func AddMatchDays(db *gorm.DB) {
	var teams []models.Country
	db.Find(&teams)

	for _, team := range teams {
		matches := []models.Match{}

		db.Where(&models.Match{AID: team.ID}).Or(&models.Match{BID: team.ID}).Find(&matches).Order("when")

		for index, match := range matches {
			if match.Day == 0 && match.Stage == models.GROUP {
				match.Day = int(index + 1)
			} else if match.Stage > models.GROUP {
				match.Day = -1
			}
			match.Assigned = true
			db.Save(&match)
		}
	}
}
