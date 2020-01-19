package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/smolvitos/idoor/internal/repository"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) ImitateUserAlreadyLoggedIn() error {
	admin, err := s.FindUserByLogin("admin")
	if err != nil {
		return err
	}
	token := s.GenerateAuthToken(admin)
	err = s.SaveTokenForAuth(admin, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) FindOneUser(id uint) (*repository.User, error) {
	var user repository.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) FindUserByToken(token string) (*repository.User, error) {
	var user repository.User
	if err := s.db.Where("token = ?", token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (Service) GenerateAuthToken(user *repository.User) string {
	userID := strconv.FormatUint(uint64(user.ID), 10)
	md5Hash := md5.Sum([]byte(userID))
	return fmt.Sprintf("%s", hex.EncodeToString(md5Hash[:]))
}

func (s *Service) SaveTokenForAuth(user *repository.User, token string) error {
	user.Token = token
	if err := s.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) FindUserByLogin(login string) (*repository.User, error) {
	var user repository.User
	if err := s.db.Where("login = ?", login).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) FindAllUsers() ([]*repository.User, error) {
	users := make([]*repository.User, 0)
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) FindMessages(from, to uint) ([]*repository.Message, error) {
	messages := make([]*repository.Message, 0)
	if err := s.db.
		Where("source = ? AND destination = ?", from, to).
		Or("source = ? AND destination = ?", to, from).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
