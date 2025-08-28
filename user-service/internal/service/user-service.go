package service

import (
	"errors"
	"time"
	"user-service/internal/domain/model"
	"user-service/internal/domain/ports"
	"user-service/pkg/hash"
	"user-service/pkg/jwt"
)

type UserService struct {
	Repo ports.UserRepository
}

func (s *UserService) Register(user *model.User) error {
	hashPassword, _ := hash.HashPassword(user.Password)
	user.Password = hashPassword
	createdAt := time.Now()
	user.CreatedAt = createdAt.Format("2006-01-02")
	return s.Repo.Save(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	comparedSucces := hash.ComparePassword(user.Password, password)
	if comparedSucces != nil {
		return "", errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(int(user.ID),user.Email)
	return token, err
}

func (s *UserService) RecoverPassword(email string) error {
	return s.Repo.RecoverPassword(email)
}

func (s *UserService) GetId(userId int) (*model.User, error) {
	return s.Repo.GetId(userId)
}

func (s *UserService) Update(email, password string) error {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return err
	}
	hashPassword, _ := hash.HashPassword(password)
	user.Password = hashPassword
	return s.Repo.Update(user)
}

func (s *UserService) Delete(email string) error {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return err
	}
	return s.Repo.Delete(user)
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.Repo.GetAll()
}
