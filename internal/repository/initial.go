package repository

import (
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
		Password: "Сюда надо что нибудь посложней, чем просто password.",
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

	if err := db.
		Create(&Message{
			Source:      test.ID,
			Destination: test.ID,
			Text:        "Ура, диалог с самим собой)",
		}).
		Create(&Message{
			Source:      user.ID,
			Destination: admin.ID,
			Text:        "Скинь ту инфу, про которую ты говорил",
		}).
		Create(&Message{
			Source:      admin.ID,
			Destination: user.ID,
			Text:        code,
		}).
		Create(&Message{
			Source:      test.ID,
			Destination: user.ID,
			Text:        "Гы, тест",
		}).
		Create(&Message{
			Source:      test.ID,
			Destination: user.ID,
			Text:        "У нас регистрация по емэйлам.",
		}).
		Create(&Message{
			Source:      user.ID,
			Destination: test.ID,
			Text:        "В смысле?",
		}).
		Create(&Message{
			Source:      test.ID,
			Destination: user.ID,
			Text:        "Кто-то особо умный тестирует систему, вбивая test@test.ch",
		}).
		Create(&Message{
			Source:      test.ID,
			Destination: user.ID,
			Text:        "Только что приходит письмо с test.ch. В переводе с французского: \"Вы ЗАЕБАЛИ!\" :)",
		}).
		Create(&Message{
			Source:      user.ID,
			Destination: test.ID,
			Text:        "%;)",
		}).
		Error; err != nil {
		return err
	}
	return nil
}
