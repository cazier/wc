package models

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model `json:"-"`

	ID       int    `gorm:"primarykey" json:"id" uri:"id"`
	Name     string `gorm:"unique" json:"name" uri:"name"`
	Group    string `json:"group" uri:"group"`
	FifaCode string `gorm:"unique" json:"fifa_code" uri:"code"`
}

type Player struct {
	gorm.Model `json:"-"`

	ID        int     `gorm:"primarykey" json:"id" uri:"id"`
	CountryID int     `json:"-"`
	Country   Country `json:"country"`

	Name     string `json:"name" uri:"name"`
	Position string `json:"position"`
	Number   int    `gorm:"default:-1" json:"number"`

	Goals  uint `gorm:"default:0" json:"goals" `
	Yellow uint `gorm:"default:0" json:"yellows"`
	Red    uint `gorm:"default:0" json:"reds"`
	Saves  int  `gorm:"default:-1" json:"saves"`
}
