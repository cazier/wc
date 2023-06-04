package load

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cazier/wc/db"
	"github.com/cazier/wc/db/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

func init() {
	db.InitSqlite(&db.SqliteDBOptions{Memory: true, LogLevel: 3})
	db.LinkTables(false)
}

var TempDir string
var TestCountryData = map[string]string{
	"name":  "Country %s",
	"code":  "C_%s",
	"group": "%s",
}
var Characters = []string{"A", "B", "C", "D", "E", "F", "G", "H"}

func TestMain(m *testing.M) {
	TempDir, _ = os.MkdirTemp("", "go_test")

	status := m.Run()

	os.RemoveAll(TempDir)
	os.Exit(status)
}

func createYaml(data interface{}, file string) string {
	var text []byte
	path := filepath.Join(TempDir, file)

	switch data := data.(type) {
	case string:
		text = []byte(data)
	default:
		text, _ = yaml.Marshal(data)
	}

	os.WriteFile(path, text, os.ModePerm)

	return path
}

func TestEmpty(t *testing.T) {
	for _, model := range []interface{}{models.Country{}, models.Player{}, models.Match{}} {
		num := int64(1)
		assert.NotZero(t, num)

		db.Database.Model(&model).Count(&num)
		assert.Zero(t, num)
	}
}

func TestTeams(t *testing.T) {
	assert := assert.New(t)
	testData := make([]map[string]string, len(Characters))

	for index, char := range Characters {
		testData[index] = make(map[string]string)
		for key, value := range TestCountryData {
			testData[index][key] = fmt.Sprintf(value, char)
		}
	}

	path := createYaml(testData, "teams.yaml")
	Teams(path)

	var num int64
	var rows []models.Country

	db.Database.Model(&rows).Count(&num)
	db.Database.Find(&rows)

	assert.Len(rows, len(testData)+2)
	assert.Equal("<A>", rows[0].FifaCode)
	assert.Equal("<B>", rows[1].FifaCode)

	for index := 2; index < int(num); index++ {
		assert.NotZero(rows[index].Name)
		assert.NotZero(rows[index].FifaCode)
		assert.NotZero(rows[index].Group)
	}
}

func TestMatches(t *testing.T) {
	assert := assert.New(t)
	testData := make([]map[string]string, len(Characters)*len(Characters))

	var counter int = 0

	for _, char_a := range Characters {
		for _, char_b := range Characters {
			if char_a == char_b {
				continue
			}

			testData[counter] = map[string]string{
				"a":     fmt.Sprintf(TestCountryData["name"], char_a),
				"b":     fmt.Sprintf(TestCountryData["name"], char_b),
				"date":  "01-Jan-01",
				"time":  fmt.Sprintf("%02d:%02d", counter/60, counter%60),
				"stage": "GROUP",
			}

			counter += 1
		}
	}

	TestTeams(t)

	testData = testData[:counter]

	path := createYaml(testData, "matches.yaml")
	Matches(path)

	var num int64
	var rows []models.Match

	db.Database.Model(&rows).Count(&num)
	db.Database.Joins("ACountry").Joins("BCountry").Find(&rows)

	assert.Len(rows, int(num))

	for index := 0; index < int(num); index++ {
		assert.NotZero(rows[index].ACountry.Name)
		assert.NotZero(rows[index].BCountry.Name)
	}

	Matches(path)
	db.Database.Model(&rows).Count(&num)
	assert.Len(rows, int(num))
}

func TestPlayers(t *testing.T) {
	assert := assert.New(t)
	testData := make([]map[string]any, len(Characters)*len(Characters)*len(Characters))

	var counter int = 0

	for _, char_a := range Characters {
		for _, char_b := range Characters {
			for _, char_c := range Characters {
				testData[counter] = map[string]any{
					"name":     fmt.Sprintf("%s %s %s", char_a, char_b, char_c),
					"country":  fmt.Sprintf(TestCountryData["name"], char_a),
					"number":   counter + 1,
					"position": fmt.Sprintf("%s %s", char_b, char_c),
				}

				counter += 1
			}
		}
	}

	TestTeams(t)

	testData = testData[:counter]

	path := createYaml(testData, "players.yaml")
	Players(path)

	var num int64
	var rows []models.Player

	db.Database.Model(&rows).Count(&num)
	db.Database.Joins("Country").Find(&rows)

	assert.Len(rows, int(num))

	for index := 0; index < int(num); index++ {
		assert.NotZero(rows[index].Name)
		assert.NotZero(rows[index].Number)
		assert.NotZero(rows[index].Position)
		assert.NotZero(rows[index].Country.Name)
	}

	Players(path)
	db.Database.Model(&rows).Count(&num)
	assert.Len(rows, int(num))
}

func TestCache(t *testing.T) {
	assert := assert.New(t)

	// Depending on the test order, this may have been filled by the TestTeam function
	cache = nil

	sql, _ := db.Database.DB()
	sql.Close()

	_, err := cacheCountries()
	assert.ErrorContains(err, "sql: database is closed")

	db.Database, _ = gorm.Open(db.Database.Dialector)
	db.LinkTables(true)

	_, err = cacheCountries()
	assert.ErrorContains(err, "cannot import match data when there are no countries in the table")

	TestTeams(t)

	output, err := cacheCountries()

	assert.NotEmpty(cache)
	assert.NotEmpty(output)
	assert.Nil(err)
}
