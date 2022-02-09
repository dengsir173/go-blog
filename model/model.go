package model

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Blog struct {
	gorm.Model
	Title   string
	Content string `gorm:"type:text"`
	Tag     string
}

type Image struct {
	gorm.Model
	Url string
}
