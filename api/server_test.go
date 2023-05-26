package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/cazier/wc/api/exceptions"
	"github.com/cazier/wc/db"
	"github.com/cazier/wc/db/load"
	"github.com/cazier/wc/db/load/utils"
	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/version"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var m Mock

func init() {
	gin.SetMode(gin.ReleaseMode)
	db.InitSqlite(&db.SqliteDBOptions{Memory: true, LogLevel: 3})
	db.LinkTables(false)

	load.Teams("../test/teams.yaml")
	load.Matches("../test/matches.yaml")
	load.Players("../test/players.yaml")

	m = Mock{
		engine:   gin.New(),
		response: *httptest.NewRecorder(),
	}
	setupRoutes(m.engine)
}

type Mock struct {
	engine   *gin.Engine
	response httptest.ResponseRecorder
}

type Response struct {
	status int
	body   string
	json   map[string]any
}

func (m *Mock) request(method, endpoint string) Response {
	var response map[string]any
	m.response = *httptest.NewRecorder()

	req, _ := http.NewRequest(method, endpoint, nil)
	m.engine.ServeHTTP(&m.response, req)

	json.Unmarshal(m.response.Body.Bytes(), &response)

	return Response{body: m.response.Body.String(), json: response, status: m.response.Code}
}

func (m *Mock) GET(endpoint string) Response {
	return m.request("GET", endpoint)
}

func (m *Mock) POST(endpoint string) Response {
	return m.request("POST", endpoint)
}

func assertException(t *testing.T, response Response, status int, exception error, messages ...string) {
	assert.Equal(t, status, response.status, messages)
	assert.Equal(t, map[string]any{"error": exceptions.Message(exception)}, response.json, messages)
}

func nilMap(m any) map[string]any {
	var output map[string]any
	data, _ := json.Marshal(m)
	json.Unmarshal(data, &output)
	return output
}

func loadPlayer() models.Player {
	var test models.Player
	players := utils.LoadPlayers("../test/players.yaml")
	player := players[rand.Intn(len(players))]

	storage, _ := json.Marshal(player)
	json.Unmarshal(storage, &test)

	return test
}

func loadCountry() models.Country {
	var test models.Country
	countries := utils.LoadTeams("../test/teams.yaml")
	country := countries[rand.Intn(len(countries))]

	storage, _ := json.Marshal(country)
	json.Unmarshal(storage, &test)
	// TODO: Use FifaCode everywhere
	test.FifaCode = country.Code

	return test
}

func TestVersion(t *testing.T) {
	response := m.GET("/version")

	assert.Equal(t, 200, response.status)
	assert.Equal(t, map[string]any{"version": version.Version}, response.json)
}

func testPlayer(t *testing.T, player models.Player, response Response) {
	assert := assert.New(t)
	data := response.json["data"].(map[string]any)
	nilPlayer := models.Player{}

	assert.Equal(http.StatusOK, response.status)

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

func testCountry(t *testing.T, country models.Country, response Response) {
	assert := assert.New(t)
	data := response.json["data"].(map[string]any)
	nilCountry := models.Country{}

	assert.Equal(http.StatusOK, response.status)

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

func testMatch(t *testing.T, response Response) {
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
	assert.Equal(http.StatusOK, response.status)

	data := response.json["data"]
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

	data := utils.LoadPlayers("../test/players.yaml")

	assert.Equal(t, 200, response.status)
	assert.Len(t, response.json["data"], len(data))

	response = Response{
		json:   map[string]any{"data": response.json["data"].([]any)[rand.Intn(len(data))]},
		status: response.status,
	}

	testPlayer(t, models.Player{}, response)
}

func TestPlayerName(t *testing.T) {
	player := loadPlayer()
	response := m.GET(fmt.Sprintf("/player/name/%s", player.Name))

	testPlayer(t, player, response)
}

func TestPlayerId(t *testing.T) {
	player := loadPlayer()
	player.ID = int(m.GET(fmt.Sprintf("/player/name/%s", player.Name)).json["data"].(map[string]any)["id"].(float64))
	response := m.GET(fmt.Sprintf("/player/id/%d", player.ID))

	testPlayer(t, player, response)
}
func TestCountries(t *testing.T) {
	response := m.GET("/country")

	data := utils.LoadTeams("../test/teams.yaml")

	assert.Equal(t, 200, response.status)
	assert.Len(t, response.json["data"], len(data))

	response = Response{
		json:   map[string]any{"data": response.json["data"].([]any)[rand.Intn(len(data))]},
		status: response.status,
	}

	testCountry(t, models.Country{}, response)
}

func TestCountryName(t *testing.T) {
	country := loadCountry()
	response := m.GET(fmt.Sprintf("/country/name/%s", country.Name))

	testCountry(t, country, response)
}
func TestCountryId(t *testing.T) {
	country := loadCountry()
	country.ID = int(m.GET(fmt.Sprintf("/country/name/%s", country.Name)).json["data"].(map[string]any)["id"].(float64))
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

	for _, player := range utils.LoadPlayers("../test/players.yaml") {
		if player.Country == country.FifaCode {
			count++
		}
	}

	assert.Len(t, response.json["data"], count)
	for _, player := range response.json["data"].([]any) {
		assert.NotContains(t, player.(map[string]any), "country")
	}
	id := int(m.GET(fmt.Sprintf("/country/name/%s", country.Name)).json["data"].(map[string]any)["id"].(float64))
	response = m.GET(fmt.Sprintf("/country/id/%d/players", id))

	assert.Len(t, response.json["data"], count)
	for _, player := range response.json["data"].([]any) {
		assert.NotContains(t, player.(map[string]any), "country")
	}
}

func TestCountryMatches(t *testing.T) {
	var count int
	country := loadCountry()
	country.ID = int(m.GET(fmt.Sprintf("/country/name/%s", country.Name)).json["data"].(map[string]any)["id"].(float64))

	response := m.GET(fmt.Sprintf("/country/id/%d/matches", country.ID))
	for index, match := range m.GET(fmt.Sprintf("/country/name/%s/matches", country.Name)).json["data"].([]any) {
		assert.Equal(t, match, response.json["data"].([]any)[index])
	}

	for _, match := range utils.LoadMatches("../test/matches.yaml") {
		if match.A == country.Name || match.B == country.Name {
			count++
		}
	}

	assert.Len(t, response.json["data"], count)
	for _, match := range response.json["data"].([]any) {
		assert.NotZero(t, match.(map[string]any)["country_a"].(map[string]any)["id"])
		assert.NotZero(t, match.(map[string]any)["country_b"].(map[string]any)["id"])
	}
}

func TestPlayerMatches(t *testing.T) {
	var count int
	player := loadPlayer()
	request := m.GET(fmt.Sprintf("/player/name/%s", player.Name)).json["data"].(map[string]any)

	player.ID = int(request["id"].(float64))
	player.Country.Name = request["country"].(map[string]any)["name"].(string)

	response := m.GET(fmt.Sprintf("/player/id/%d/matches", player.ID))
	for index, match := range m.GET(fmt.Sprintf("/player/name/%s/matches", player.Name)).json["data"].([]any) {
		assert.Equal(t, match, response.json["data"].([]any)[index])
	}

	for _, match := range utils.LoadMatches("../test/matches.yaml") {
		if match.A == player.Country.Name || match.B == player.Country.Name {
			count++
		}
	}

	assert.Len(t, response.json["data"], count)
	for _, match := range response.json["data"].([]any) {
		assert.NotZero(t, match.(map[string]any)["country_a"].(map[string]any)["id"])
		assert.NotZero(t, match.(map[string]any)["country_b"].(map[string]any)["id"])
	}
}

func TestMatches(t *testing.T) {
	response := m.GET("/match")

	data := utils.LoadMatches("../test/matches.yaml")

	assert.Equal(t, 200, response.status)
	assert.Len(t, response.json["data"], len(data))

	response = Response{
		json:   map[string]any{"data": response.json["data"].([]any)[rand.Intn(len(data))]},
		status: response.status,
	}

	testMatch(t, response)
}

func TestMatchId(t *testing.T) {
	id := rand.Intn(len(utils.LoadMatches("../test/matches.yaml")))
	response := m.GET(fmt.Sprintf("/match/id/%d", id))

	testMatch(t, response)
}

func TestMatchDay(t *testing.T) {
	// TODO Remove hardcoded values
	response := m.GET(fmt.Sprintf("/match/day/%d", rand.Intn(2)+1))
	assert.Len(t, response.json["data"].([]any), 16)
	testMatch(t, response)
}
func TestMatchGroup(t *testing.T) {
	response := m.GET(fmt.Sprintf("/match/group/%s", []string{"A", "B", "C", "D", "E", "F", "G", "H"}[rand.Intn(8)]))
	// TODO Remove hardcoded values
	assert.Len(t, response.json["data"].([]any), 6)
	testMatch(t, response)
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
