package model

import (
	"github.com/jinzhu/gorm"
)

type Plugin struct {
	gorm.Model

	Name     string
	Intro    string
	Thumb    string
	Status   int8
	Filename string
}
