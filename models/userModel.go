package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Userid int64 `gorm:"unique"`
}
