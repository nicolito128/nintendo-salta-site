package models

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	Name   string `gorm:"varchar(50); not null" json:"name"`
	Token  string `gorm:"varchar(250); not null; unique" json:"token"`
	Expire int64  `gorm:"not null" json:"expire"`
}
