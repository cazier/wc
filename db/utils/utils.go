package utils

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Team struct {
	Name  string
	Code  string
	Group string
}

type Match struct {
	A     string
	B     string
	Group string
	Date  time.Time
}

func (m *Match) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var base struct {
		A     string
		B     string
		Group string
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
	m.Group = base.Group
	m.Date = time.Date(dd.Year(), dd.Month(), dd.Day(), tt.Hour(), tt.Minute(), 0, 0, time.UTC)

	return nil
}

func Import() ([]Team, []Match) {
	f, _ := os.ReadFile("teams.yaml")
	t := []Team{}
	_ = yaml.Unmarshal(f, &t)

	f2, _ := os.ReadFile("matches.yaml")
	m := []Match{}
	_ = yaml.Unmarshal(f2, &m)

	return t, m
}
