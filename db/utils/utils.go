package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Team struct {
	Name  string
	Code  string
	Group string
}

type Player struct {
	Name      string
	Number    int
	Postition string
}

type Match struct {
	A     string
	B     string
	Stage Stage
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

	err := unmarshal(&base)

	if err != nil {
		panic(err)
	}

	dd, err := time.Parse("02-Jan-06", base.Date)

	if err != nil {
		panic(err)
	}

	tt, err := time.Parse("15:04", base.Time)

	if err != nil {
		panic(err)
	}

	m.A = base.A
	m.B = base.B
	m.Stage = UnmarshalText(base.Stage)
	m.Date = time.Date(dd.Year(), dd.Month(), dd.Day(), tt.Hour(), tt.Minute(), 0, 0, time.UTC)

	return nil
}

type Stage uint

const (
	GROUP Stage = iota
	ROUND_OF_SIXTEEN
	QUARTERFINALS
	SEMIFINALS
	THIRD_PLACE
	FINAL
)

func UnmarshalText(s string) Stage {
	switch s {
	case "GROUP":
		return GROUP
	case "ROUND_OF_SIXTEEN":
		return ROUND_OF_SIXTEEN
	case "QUARTERFINALS":
		return QUARTERFINALS
	case "SEMIFINALS":
		return SEMIFINALS
	case "THIRD_PLACE":
		return THIRD_PLACE
	case "FINAL":
		return FINAL
	}
	log.Fatalf("Could not parse stage value: %s", s)
	return 6
}

type HasName interface {
	Name() string
}

func load(path string, i interface{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read the yaml file %s because: %w", path, err))
	}

	err = yaml.Unmarshal(data, i)

	if err != nil {
		panic(fmt.Errorf("could not parse the yaml file %s because: %w", path, err))
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

func LoadPlayers(path string) map[string][]Player {
	players := map[string][]Player{}
	load(path, &players)

	return players
}
