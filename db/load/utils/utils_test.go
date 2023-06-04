package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cazier/wc/db/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var TempDir string

func TestMain(m *testing.M) {
	TempDir, _ = os.MkdirTemp("", "go_test")

	status := m.Run()

	os.RemoveAll(TempDir)
	os.Exit(status)
}

func Test_load(t *testing.T) {
	path := filepath.Join(TempDir, "loadfile")

	assert.PanicsWithError(
		t, fmt.Sprintf("could not read the yaml file %s because the file doesn't exist", path), func() {
			load(path, "")
		})

	os.WriteFile(path, []byte("-"), os.ModePerm)
	assert.PanicsWithError(
		t, fmt.Sprintf("could not parse the yaml file: %s", path), func() {
			load(path, "")
		})
}

func TestLoadTeams(t *testing.T) {
	testData := `- name: Country A
  code: C_A
  group: A
- name: Country B
  code: C_B
  group: B
- name: Country C
  code: C_C
  group: C
- name: Country D
  code: C_D
  group: D
`
	os.WriteFile(filepath.Join(TempDir, "teams.yaml"), []byte(testData), os.ModePerm)

	data := LoadTeams(filepath.Join(TempDir, "teams.yaml"))
	letters := []string{"A", "B", "C", "D"}

	for index, team := range data {
		assert.IsType(t, Team{}, team)
		assert.Equal(t, fmt.Sprintf("Country %s", letters[index]), team.Name)
		assert.Equal(t, fmt.Sprintf("C_%s", letters[index]), team.Code)
		assert.Equal(t, letters[index], team.Group)
	}
}

func TestLoadMatches(t *testing.T) {
	testData := `- a: Country A_1
  b: Country A_2
  date: 01-Jan-01
  stage: GROUP
  time: '1:00'
- a: Country B_1
  b: Country B_2
  date: 02-Feb-02
  stage: ROUND_OF_SIXTEEN
  time: '2:00'
- a: Country C_1
  b: Country C_2
  date: 03-Mar-03
  stage: QUARTERFINALS
  time: '3:00'
- a: Country D_1
  b: Country D_2
  date: 04-Apr-04
  stage: SEMIFINALS
  time: '4:00'
- a: Country E_1
  b: Country E_2
  date: 05-May-05
  stage: THIRD_PLACE
  time: '5:00'
- a: Country F_1
  b: Country F_2
  date: 06-Jun-06
  stage: FINAL
  time: '6:00'
`
	os.WriteFile(filepath.Join(TempDir, "matches.yaml"), []byte(testData), os.ModePerm)

	data := LoadMatches(filepath.Join(TempDir, "matches.yaml"))
	letters := []string{"A", "B", "C", "D", "E", "F"}

	for index, match := range data {
		assert.IsType(t, Match{}, match)
		assert.Equal(t, fmt.Sprintf("Country %s_1", letters[index]), match.A)
		assert.Equal(t, fmt.Sprintf("Country %s_2", letters[index]), match.B)
		assert.Equal(t, time.Date(index+2001, time.Month(index+1), index+1, index+1, 0, 0, 0, time.UTC), match.Date)
		assert.Equal(t, models.Stage(index), match.Stage)
	}
}

func TestMatchUnmarshalBad(t *testing.T) {
	date := "January 01, 2001"
	testData := fmt.Sprintf("a: Country A_1\nb: Country A_2\ndate: %s\nstage: GROUP\ntime: '1:00'", date)

	assert.PanicsWithError(t, fmt.Sprintf("could not parse the date from the yaml file: `%s`", date),
		func() {
			yaml.Unmarshal([]byte(testData), &Match{})
		})

	time := "1 AM"
	testData = fmt.Sprintf("a: Country A_1\nb: Country A_2\ndate: 01-Jan-01\nstage: GROUP\ntime: '%s'", time)

	assert.PanicsWithError(t, fmt.Sprintf("could not parse the time from the yaml file: `%s`", time),
		func() {
			yaml.Unmarshal([]byte(testData), &Match{})
		})

	assert.Panics(t, func() { UnmarshalText("INVALID_STAGE") })
}

func TestLoadPlayers(t *testing.T) {
	testData := `- name: "First Middle Last"
  country: ABC
  number: 1
  position: GK
`
	os.WriteFile(filepath.Join(TempDir, "players.yaml"), []byte(testData), os.ModePerm)

	data := LoadPlayers(filepath.Join(TempDir, "players.yaml"))

	assert.IsType(t, []Player{}, data)
	assert.EqualValues(t, []Player{{Name: "First Middle Last", Country: "ABC", Number: 1, Position: "GK"}}, data)
}
