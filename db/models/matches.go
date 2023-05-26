package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model `json:"-"`

	ID     int  `gorm:"primarykey" json:"id"`
	Day    int  `gorm:"default:0"  json:"match_day" uri:"day"`
	Played bool `gorm:"default:false" json:"played"`

	AID      int     `json:"-"`
	BID      int     `json:"-"`
	ACountry Country `gorm:"foreignKey:AID" json:"country_a"`
	BCountry Country `gorm:"foreignKey:BID" json:"country_b"`

	Stage Stage `json:"-"`

	When     time.Time `json:"when"`
	Assigned bool      `gorm:"default:false" json:"-"`

	// AResult MatchResult `gorm:"foreignKey:ID"`
	// BResult MatchResult `gorm:"foreignKey:ID"`
}

type MatchResult struct {
	gorm.Model `json:"-"`

	ID      int `gorm:"primarykey" json:"id"`
	MatchID int

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
