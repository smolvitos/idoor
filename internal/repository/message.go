package repository

import "github.com/jinzhu/gorm"

type Message struct {
	gorm.Model
	Source      uint
	Destination uint
	Text        string `gorm:"size:25000"`
}
