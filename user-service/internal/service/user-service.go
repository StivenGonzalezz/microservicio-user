package service

import (
	"errors"
	"fmt"
	"time"
	"user-service/internal/domain/model"
	"user-service/internal/domain/ports"
	"user-service/pkg/hash"
	"user-service/pkg/jwt"
	"user-service/pkg/rabbitmq"
)

type UserService struct {
	Repo      ports.UserRepository
	Publisher *rabbitmq.Publisher
}

func (s *UserService) Register(user *model.User) error {
	hashPassword, _ := hash.HashPassword(user.Password)
	user.Password = hashPassword
	createdAt := time.Now()
	user.CreatedAt = createdAt.Format("2006-01-02")

	err := s.Repo.Save(user)
	if err != nil {
		return err
	}

	// publico evento en RabbitMQ
	event := map[string]interface{}{
		"action": "user.registered",
		"user": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"lastName": user.LastName,
			"email":    user.Email,
			"phone":    user.PhoneNumber,
		},
		"timestamp": time.Now(),
	}

	err = s.Publisher.Publish("user.registered", event)
	if err != nil {
		return err
	}

	return nil
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

	token, err := jwt.GenerateToken(int(user.ID), user.Email)
	if err != nil {
		return "", err
	}

	// publico evento en RabbitMQ
	event := map[string]interface{}{
		"action": "user.login",
		"user": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"lastName": user.LastName,
			"email":    user.Email,
			"phone":    user.PhoneNumber,
		},
		"timestamp": time.Now(),
	}

	err = s.Publisher.Publish("user.login", event)
	if err != nil {
		return "", err
	}

	return token, err
}

func (s *UserService) RecoverPassword(email string) (string, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	// De momento solo usamos el ID en la URL (más adelante podemos usar un token seguro)
	recoveryURL := fmt.Sprintf("http://localhost:8080/user/password/%d", user.ID)

	event := map[string]interface{}{
		"action": "user.recovery.link",
		"user": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"lastName": user.LastName,
			"email":    user.Email,
			"phone":    user.PhoneNumber,
		},
		"timestamp": time.Now(),
	}

	err = s.Publisher.Publish("user.recovery.link", event)
	if err != nil {
		return "", err
	}

	return recoveryURL, nil
}

func (s *UserService) UpdatePassword(id uint, password string) error {
	user, err := s.Repo.GetId(int(id))
	if err != nil {
		return err
	}

	hashPassword, _ := hash.HashPassword(password)
	user.Password = hashPassword

	err = s.Repo.Update(user)
	if err != nil {
		return err
	}

	event := map[string]interface{}{
		"action": "user.password.updated",
		"user": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"lastName": user.LastName,
			"email":    user.Email,
			"phone":    user.PhoneNumber,
		},
		"timestamp": time.Now(),
	}

	err = s.Publisher.Publish("user.password.updated", event)
	if err != nil {
		return err
	}

	return nil
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
	if comparedSucces != nil {
		return errors.New("invalid password")
	}
	return s.Repo.Delete(user)
}

func (s *UserService) GetByName(nameOrEmail string) ([]model.User, error) {
	return s.Repo.GetByName(nameOrEmail)
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) GetUsersWithPagination(name string, page, limit int, sort string) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := s.Repo.GetWithPagination(name, page, limit, sort)
	if err != nil {
		return nil, fmt.Errorf("error getting users with pagination: %w", err)
	}

	response := map[string]interface{}{
		"total":      result.Total,
		"page":       result.Page,
		"limit":      result.Limit,
		"totalPages": result.TotalPages,
		"data":       result.Data,
	}

	return response, nil
}