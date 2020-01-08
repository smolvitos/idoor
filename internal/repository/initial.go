package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//CreateInitialData заполняет базу тестовыми данными
func CreateInitialData(db *gorm.DB, code string) error {
	admin := User{
		Login:    "admin",
		Password: "Pa$$vv0rcl",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	user := User{
		Login:    "user",
		Password: "password",
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	test := User{
		Login:    "test",
		Password: "test",
	}
	if err := db.Create(&test).Error; err != nil {
		return err
	}

	if err := db.Create(&Message{
		Source:      admin.ID,
		Destination: user.ID,
		Text:        fmt.Sprintf("Код запуска ракет: %s", code),
	}).Create(&Message{
		Source:      test.ID,
		Destination: user.ID,
		Text:        "Гы, тест",
	}).Error; err != nil {
		return err
	}
	return nil
}
