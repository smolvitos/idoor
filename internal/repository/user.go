package repository

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Login    string `gorm:"unique_index"`
	Password string
	Token    string
}
