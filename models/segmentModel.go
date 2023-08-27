package models

import "gorm.io/gorm"

type Segment struct {
	gorm.Model
	Slug string `gorm:"unique"`
}
