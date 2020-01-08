package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/dilap54/voronov_idor/internal/repository"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) FindOneUser(id uint) (*repository.User, error) {
	var user repository.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) FindUserByToken(token string) (*repository.User, error) {
	tokenArr := strings.Split(token, ":")
	if len(tokenArr) != 2 {
		return nil, fmt.Errorf("token should consist of 2 parts divided by ':', got: %s", tokenArr)
	}
	userID, err := strconv.ParseUint(tokenArr[0], 10, 64)
	if err != nil {
		return nil, err
	}
	var user repository.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	if s.GenerateAuthToken(&user) == token {
		return &user, nil
	}
	return nil, nil
}

func (Service) GenerateAuthToken(user *repository.User) string {
	userID := strconv.FormatUint(uint64(user.ID), 10)
	md5Hash := md5.Sum([]byte(userID))
	return fmt.Sprintf("%s:%s", userID, hex.EncodeToString(md5Hash[:]))
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
