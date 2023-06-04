package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/cazier/wc/db/models"
	"gopkg.in/yaml.v3"
)

type Team struct {
	Name  string
	Code  string
	Group string
}

type Player struct {
	Name     string
	Country  string
	Number   int
	Position string
}

type Match struct {
	A     string
	B     string
	Stage models.Stage
	Date  time.Time
}

func (m *Match) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var base struct {
		A     string
		B     string
		Stage string
		Date  string
		Time  string
	}

	var tt time.Time
	var dd time.Time
	var err error

	unmarshal(&base)

	if dd, err = time.Parse("02-Jan-06", base.Date); err != nil {
		log.Printf("error: %s", err.Error())
		panic(fmt.Errorf("could not parse the date from the yaml file: `%s`", base.Date))
	}

	if tt, err = time.Parse("15:04", base.Time); err != nil {
		log.Printf("error: %s", err.Error())
		panic(fmt.Errorf("could not parse the time from the yaml file: `%s`", base.Time))
	}

	m.A = base.A
	m.B = base.B
	m.Stage = UnmarshalText(base.Stage)
	m.Date = time.Date(dd.Year(), dd.Month(), dd.Day(), tt.Hour(), tt.Minute(), 0, 0, time.UTC)

	return nil
}

func UnmarshalText(s string) models.Stage {
	switch s {
	case "GROUP":
		return models.GROUP
	case "ROUND_OF_SIXTEEN":
		return models.ROUND_OF_SIXTEEN
	case "QUARTERFINALS":
		return models.QUARTERFINALS
	case "SEMIFINALS":
		return models.SEMIFINALS
	case "THIRD_PLACE":
		return models.THIRD_PLACE
	case "FINAL":
		return models.FINAL
	}
	panic(fmt.Errorf("could not parse stage value: %s", s))
}

func load(path string, i interface{}) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		panic(fmt.Errorf("could not read the yaml file %s because the file doesn't exist", path))
	}

	err = yaml.Unmarshal(data, i)

	if err != nil {
		log.Printf("error: %s", err.Error())
		panic(fmt.Errorf("could not parse the yaml file: %s", path))
	}
}

func LoadTeams(path string) []Team {
	team := []Team{}
	load(path, &team)

	return team
}

func LoadMatches(path string) []Match {
	match := []Match{}
	load(path, &match)

	return match
}

func LoadPlayers(path string) []Player {
	players := []Player{}
	load(path, &players)

	return players
}
