package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cazier/wc/db/load"
	"github.com/cazier/wc/db/load/utils"
	"github.com/cazier/wc/db/models"
	test "github.com/cazier/wc/testing"
	"github.com/cazier/wc/version"
	"github.com/cazier/wc/web/exceptions"
	"github.com/stretchr/testify/assert"
)

var m test.Mock

func init() {
	m = test.NewMock(
		&test.MockOptions{
			Callback: Init,
			Models: []any{
				&models.Country{},
				&models.Player{},
				&models.Match{},
			},
		},
	)

	m.BasePath = group.BasePath()

	load.Init(m.Database)

	load.Teams(test.Path("teams.yaml"))
	load.Matches(test.Path("matches.yaml"))
	load.Players(test.Path("players.yaml"))
}

func assertException(t *testing.T, response test.Response, status int, exception error, messages ...string) {
	assert.Equal(t, status, response.Status, messages)
	assert.Equal(t, map[string]any{"error": exceptions.Message(exception)}, response.Json, messages)
}

func nilMap(m any) map[string]any {
	var output map[string]any
	data, _ := json.Marshal(m)
	json.Unmarshal(data, &output)
	return output
}

func loadPlayer() models.Player {
	var dest models.Player
	players := utils.LoadPlayers(test.Path("players.yaml"))
	player := players[rand.Intn(len(players))]

	storage, _ := json.Marshal(player)
	json.Unmarshal(storage, &dest)

	return dest
}

func loadCountry() models.Country {
	var dest models.Country
	countries := utils.LoadTeams(test.Path("teams.yaml"))
	country := countries[rand.Intn(len(countries))]

	storage, _ := json.Marshal(country)
	json.Unmarshal(storage, &dest)
	// TODO: Use FifaCode everywhere
	dest.FifaCode = country.Code

	return dest
}

func TestVersion(t *testing.T) {
	response := m.GET("/version")

	assert.Equal(t, 200, response.Status)
	assert.Equal(t, map[string]any{"version": version.Version}, response.Json)
}

func testPlayer(t *testing.T, player models.Player, response test.Response) {
	assert := assert.New(t)
	data := response.Json["data"].(map[string]any)
	nilPlayer := models.Player{}

	assert.Equal(http.StatusOK, response.Status)

	if player.ID != 0 {
		assert.EqualValues(player.ID, data["id"])
	} else if player != nilPlayer {
		assert.Equal(player.Name, data["name"])
		assert.Equal(player.Position, data["position"])
		assert.EqualValues(player.Number, data["number"])
		assert.NotEqual(nilMap(models.Country{}), data["country"])
	} else {
		assert.NotEqualValues(nilPlayer.ID, data["id"])
		assert.NotEqualValues(nilPlayer.Name, data["name"])
		assert.NotEqualValues(nilPlayer.CountryID, data["country"].(map[string]any)["id"])
	}
}

func testCountry(t *testing.T, country models.Country, response test.Response) {
	assert := assert.New(t)
	data := response.Json["data"].(map[string]any)
	nilCountry := models.Country{}

	assert.Equal(http.StatusOK, response.Status)

	if country.ID != 0 {
		assert.EqualValues(country.ID, data["id"])
	} else if country != nilCountry {
		assert.Equal(country.Name, data["name"])
		assert.Equal(country.Group, data["group"])
		assert.EqualValues(country.FifaCode, data["fifa_code"])
	} else {
		assert.NotEqualValues(nilCountry.ID, data["id"])
		assert.NotEqualValues(nilCountry.Name, data["name"])
		assert.NotEqualValues(nilCountry.Group, data["group"])
		assert.NotEqualValues(nilCountry.FifaCode, data["fifa_code"])
	}
}

func testMatch(t *testing.T, response test.Response) {
	assert := assert.New(t)
	test := func(val map[string]any) {
		nilMatch := models.Match{}

		assert.Contains(val, "played")
		assert.NotEqualValues(nilMatch.ID, val["id"])
		assert.NotEqualValues(nilMatch.Day, val["match_day"])
		assert.NotEqualValues(nilMatch.When, val["when"])

		assert.NotZero(val["country_a"].(map[string]any)["id"])
		assert.NotZero(val["country_b"].(map[string]any)["id"])
	}
	assert.Equal(http.StatusOK, response.Status)

	data := response.Json["data"]
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		test(data.(map[string]any))
	case reflect.Slice:
		for _, val := range data.([]any) {
			test(val.(map[string]any))
		}
	default:
		return
	}
}

func TestPlayers(t *testing.T) {
	response := m.GET("/player")

	data := utils.LoadPlayers(test.Path("players.yaml"))

	assert.Equal(t, 200, response.Status)
	assert.Len(t, response.Json["data"], len(data))

	response = test.Response{
		Json:   map[string]any{"data": response.Json["data"].([]any)[rand.Intn(len(data))]},
		Status: response.Status,
	}

	testPlayer(t, models.Player{}, response)
}

func TestPlayerName(t *testing.T) {
	player := loadPlayer()
	response := m.GET(fmt.Sprintf("/player/name/%s", player.Name))
	lower := m.GET(fmt.Sprintf("/player/name/%s", strings.ToLower(player.Name)))

	testPlayer(t, player, response)

	assert.EqualValues(t, response.Json, lower.Json)
}

func TestPlayerId(t *testing.T) {
	player := loadPlayer()
	player.ID = int(m.GET(fmt.Sprintf("/player/name/%s", player.Name)).Json["data"].(map[string]any)["id"].(float64))
	response := m.GET(fmt.Sprintf("/player/id/%d", player.ID))

	testPlayer(t, player, response)
}
func TestCountries(t *testing.T) {
	response := m.GET("/country")

	data := utils.LoadTeams(test.Path("teams.yaml"))

	assert.Equal(t, 200, response.Status)
	assert.Len(t, response.Json["data"], len(data))

	response = test.Response{
		Json:   map[string]any{"data": response.Json["data"].([]any)[rand.Intn(len(data))]},
		Status: response.Status,
	}

	testCountry(t, models.Country{}, response)
}

func TestCountryName(t *testing.T) {
	country := loadCountry()
	response := m.GET(fmt.Sprintf("/country/name/%s", country.Name))
	lower := m.GET(fmt.Sprintf("/country/name/%s", strings.ToLower(country.Name)))
	upper := m.GET(fmt.Sprintf("/country/name/%s", strings.ToUpper(country.Name)))

	testCountry(t, country, response)

	assert.EqualValues(t, response.Json, lower.Json)
	assert.EqualValues(t, response.Json, upper.Json)
}
func TestCountryId(t *testing.T) {
	country := loadCountry()
	country.ID = int(m.GET(fmt.Sprintf("/country/name/%s", country.Name)).Json["data"].(map[string]any)["id"].(float64))
	response := m.GET(fmt.Sprintf("/country/id/%d", country.ID))

	testCountry(t, country, response)
}

func TestCountryCode(t *testing.T) {
	country := loadCountry()
	response := m.GET(fmt.Sprintf("/country/code/%s", country.FifaCode))

	testCountry(t, country, response)
}

func TestCountryPlayers(t *testing.T) {
	var count int
	country := loadCountry()
	response := m.GET(fmt.Sprintf("/country/name/%s/players", country.Name))

	for _, player := range utils.LoadPlayers(test.Path("players.yaml")) {
		if player.Country == country.FifaCode {
			count++
		}
	}

	assert.Len(t, response.Json["data"], count)
	for _, player := range response.Json["data"].([]any) {
		assert.NotContains(t, player.(map[string]any), "country")
	}
	id := int(m.GET(fmt.Sprintf("/country/name/%s", country.Name)).Json["data"].(map[string]any)["id"].(float64))
	response = m.GET(fmt.Sprintf("/country/id/%d/players", id))

	assert.Len(t, response.Json["data"], count)
	for _, player := range response.Json["data"].([]any) {
		assert.NotContains(t, player.(map[string]any), "country")
	}
}

func TestCountryMatches(t *testing.T) {
	var count int
	country := loadCountry()
	country.ID = int(m.GET(fmt.Sprintf("/country/name/%s", country.Name)).Json["data"].(map[string]any)["id"].(float64))

	response := m.GET(fmt.Sprintf("/country/id/%d/matches", country.ID))
	assert.Equal(
		t,
		m.GET(fmt.Sprintf("/country/name/%s/matches", country.Name)).Json["data"].([]any),
		response.Json["data"].([]any),
	)

	for _, match := range utils.LoadMatches(test.Path("matches.yaml")) {
		if match.A == country.Name || match.B == country.Name {
			count++
		}
	}

	dates := make([]time.Time, count)
	for index, match := range response.Json["data"].([]any) {
		dates[index], _ = time.Parse(time.RFC3339, match.(map[string]any)["when"].(string))
	}
	assert.IsNonDecreasing(t, dates)

	assert.Len(t, response.Json["data"], count)
	for _, match := range response.Json["data"].([]any) {
		assert.NotZero(t, match.(map[string]any)["country_a"].(map[string]any)["id"])
		assert.NotZero(t, match.(map[string]any)["country_b"].(map[string]any)["id"])
	}
}

func TestPlayerMatches(t *testing.T) {
	var count int
	player := loadPlayer()
	request := m.GET(fmt.Sprintf("/player/name/%s", player.Name)).Json["data"].(map[string]any)

	player.ID = int(request["id"].(float64))
	player.Country.Name = request["country"].(map[string]any)["name"].(string)

	response := m.GET(fmt.Sprintf("/player/id/%d/matches", player.ID))
	assert.Equal(
		t,
		m.GET(fmt.Sprintf("/player/name/%s/matches", player.Name)).Json["data"].([]any),
		response.Json["data"].([]any),
	)

	for _, match := range utils.LoadMatches(test.Path("matches.yaml")) {
		if match.A == player.Country.Name || match.B == player.Country.Name {
			count++
		}
	}

	dates := make([]time.Time, count)
	for index, match := range response.Json["data"].([]any) {
		dates[index], _ = time.Parse(time.RFC3339, match.(map[string]any)["when"].(string))
	}
	assert.IsNonDecreasing(t, dates)

	assert.Len(t, response.Json["data"], count)
	for _, match := range response.Json["data"].([]any) {
		assert.NotZero(t, match.(map[string]any)["country_a"].(map[string]any)["id"])
		assert.NotZero(t, match.(map[string]any)["country_b"].(map[string]any)["id"])
	}
}

func TestMatches(t *testing.T) {
	response := m.GET("/match")

	data := utils.LoadMatches(test.Path("matches.yaml"))

	assert.Equal(t, 200, response.Status)
	assert.Len(t, response.Json["data"], len(data))

	dates := make([]time.Time, len(data))
	for index, match := range response.Json["data"].([]any) {
		dates[index], _ = time.Parse(time.RFC3339, match.(map[string]any)["when"].(string))
	}
	assert.IsNonDecreasing(t, dates)

	response = test.Response{
		Json:   map[string]any{"data": response.Json["data"].([]any)[rand.Intn(len(data))]},
		Status: response.Status,
	}

	testMatch(t, response)
}

func TestMatchId(t *testing.T) {
	id := rand.Intn(len(utils.LoadMatches(test.Path("matches.yaml"))))
	response := m.GET(fmt.Sprintf("/match/id/%d", id))

	testMatch(t, response)
}

func TestMatchDay(t *testing.T) {
	// TODO Remove hardcoded values
	response := m.GET(fmt.Sprintf("/match/day/%d", rand.Intn(2)+1))
	assert.Len(t, response.Json["data"].([]any), 16)
	testMatch(t, response)
}
func TestMatchGroup(t *testing.T) {
	group := []string{"A", "B", "C", "D", "E", "F", "G", "H"}[rand.Intn(8)]

	response := m.GET(fmt.Sprintf("/match/group/%s", group))
	// TODO Remove hardcoded values
	assert.Len(t, response.Json["data"].([]any), 6)
	testMatch(t, response)

	lower := m.GET(fmt.Sprintf("/match/group/%s", strings.ToLower(group)))
	upper := m.GET(fmt.Sprintf("/match/group/%s", strings.ToUpper(group)))

	assert.EqualValues(t, response.Json, lower.Json)
	assert.EqualValues(t, response.Json, upper.Json)
}

func TestNameBad(t *testing.T) {
	tests := map[string][]string{"player": {"matches"}, "country": {"players", "matches"}}

	for endpoint, result := range tests {
		response := m.GET(fmt.Sprintf("/%s/name/notarealname", endpoint))
		assertException(t, response, http.StatusBadRequest, &exceptions.NoResultsFoundError{})

		for _, res := range result {
			response = m.GET(fmt.Sprintf("/%s/name/notarealname/%s", endpoint, res))
			assertException(t, response, http.StatusBadRequest, &exceptions.NoResultsFoundError{})
			response = m.GET(fmt.Sprintf("/%s/name//%s", endpoint, res))
			assertException(t, response, http.StatusUnprocessableEntity, &exceptions.InvalidValueError{})
		}
	}
}

func TestIdBad(t *testing.T) {
	tests := map[string][]string{"player": {"matches"}, "country": {"players", "matches"}, "match": {}}

	for endpoint, result := range tests {
		response := m.GET(fmt.Sprintf("/%s/id/0", endpoint))
		assertException(t, response, http.StatusBadRequest, &exceptions.NoResultsFoundError{})

		response = m.GET(fmt.Sprintf("/%s/id/invalidtype", endpoint))
		assertException(t, response, http.StatusUnprocessableEntity, &strconv.NumError{})

		for _, res := range result {
			response = m.GET(fmt.Sprintf("/%s/id/invalidtype/%s", endpoint, res))
			assertException(t, response, http.StatusUnprocessableEntity, &strconv.NumError{})
			response = m.GET(fmt.Sprintf("/%s/id/999999/%s", endpoint, res))
			assertException(t, response, http.StatusBadRequest, &exceptions.NoResultsFoundError{})
		}
	}
}
