package models

import (
	"time"

	"github.com/cazier/wc/db/utils"
	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	Day    uint `gorm:"default:0"`
	Played bool `gorm:"default:false"`

	AID   uint
	BID   uint
	Stage utils.Stage

	When     time.Time
	Assigned bool `gorm:"default:false"`

	AResult MatchResult `gorm:"foreignKey:ID"`
	BResult MatchResult `gorm:"foreignKey:ID"`
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
