package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model

	Name     string `gorm:"unique"`
	Group    string
	FifaCode string `gorm:"unique"`

	Matches []Match  `gorm:"foreignKey:ID"`
	Players []Player `gorm:"foreignKey:ID"`
}

type Player struct {
	gorm.Model
	CountryID uint

	Name     string
	Position string
	Number   int `gorm:"default:-1"`

	Goals  uint `gorm:"default:0"`
	Yellow uint `gorm:"default:0"`
	Red    uint `gorm:"default:0"`
	Saves  int  `gorm:"default:-1"`
}
