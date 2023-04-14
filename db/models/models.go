package models

import (
	"time"

	"gorm.io/gorm"
)

type Country struct {
	gorm.Model

	Name     string
	Group    string
	FifaCode string

	Matches []Match `gorm:"foreignKey:ID"`
}

type Match struct {
	gorm.Model

	Day    uint
	Played bool

	AID   uint
	BID   uint
	Group string

	When time.Time

	// AResult MatchResult
	// BResult MatchResult
}

type MatchResult struct {
	Yellow       uint
	Red          uint
	GoalsFor     uint
	GoalsAgainst uint
	Points       uint
}
