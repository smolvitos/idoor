package repository

import (
	"github.com/jinzhu/gorm"
	//Импорт sqlite, потому что он нужен для gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "test.sqlite")
	if err != nil {
		return nil, err
	}
	//db.LogMode(true)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := deleteOldTables(db); err != nil {
		return err
	}
	if err := createNewTables(db); err != nil {
		return err
	}
	return nil
}

func deleteOldTables(db *gorm.DB) error {
	if err := db.DropTableIfExists(Message{}).Error; err != nil {
		return err
	}
	if err := db.DropTableIfExists(User{}).Error; err != nil {
		return err
	}
	return nil
}

func createNewTables(db *gorm.DB) error {
	if err := db.CreateTable(User{}).Error; err != nil {
		return err
	}
	if err := db.CreateTable(Message{}).Error; err != nil {
		return err
	}
	return nil
}
