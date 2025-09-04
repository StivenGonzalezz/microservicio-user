package service

import (
	"errors"
	"fmt"
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

func (s *UserService) RecoverPassword(email string) (string, error) {
    user, err := s.Repo.GetByEmail(email)
    if err != nil {
        return "", err
    }

    // De momento solo usamos el ID en la URL (más adelante podemos usar un token seguro)
    recoveryURL := fmt.Sprintf("http://localhost:8080/user/password/%d", user.ID)

    return recoveryURL, nil
}

func (s *UserService) UpdatePassword(id uint, password string) error {
    user, err := s.Repo.GetId(int(id))
    if err != nil {
        return err
    }

    hashPassword, _ := hash.HashPassword(password)
    user.Password = hashPassword

    return s.Repo.Update(user)
}



func (s *UserService) GetId(userId int) (*model.User, error) {
	return s.Repo.GetId(userId)
}

func (s *UserService) Update(user *model.User) error {
	userdb, err := s.Repo.GetId(int(user.ID))
	if err != nil {
		return err
	}

	comparedSucces := hash.ComparePassword(userdb.Password, user.Password)
	if comparedSucces != nil {
		return errors.New("invalid password")
	}
	user.Password = userdb.Password

    return s.Repo.Update(user)
}


func (s *UserService) Delete(userId int, email string, password string) error {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user.Email != email {
		return errors.New("invalid email")
	}
	comparedSucces := hash.ComparePassword(user.Password, password)
	if comparedSucces != nil{
		return errors.New("invalid password")
	}
	return s.Repo.Delete(user)
}

func (s *UserService) GetAll(nameOrEmail string) ([]model.User, error) {
	return s.Repo.GetAll(nameOrEmail)
}
