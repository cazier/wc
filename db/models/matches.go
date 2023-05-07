package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	Day    uint `gorm:"default:0"`
	Played bool `gorm:"default:false"`

	AID   uint
	BID   uint
	Stage Stage

	When     time.Time
	Assigned bool `gorm:"default:false"`

	// AResult MatchResult `gorm:"foreignKey:ID"`
	// BResult MatchResult `gorm:"foreignKey:ID"`
}

type MatchResult struct {
	gorm.Model

	MatchID uint

	Yellow       uint
	Red          uint
	GoalsFor     uint
	GoalsAgainst uint
	Points       uint
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
