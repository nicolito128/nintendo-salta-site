package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name  string `gorm:"varchar(100); unique; not null" json:"name"`
	Score int    `gorm:"default:0" json:"score"`
}

type APIUser struct {
	Name string `json:"name"`
}
