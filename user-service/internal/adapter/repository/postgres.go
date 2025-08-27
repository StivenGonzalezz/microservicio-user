package repository

import (
	"user-service/internal/domain/model"
    "user-service/internal/domain/ports"
	"user-service/pkg/hash"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo() ports.UserRepository {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la base de datos: " + err.Error())
	}
	
	db.AutoMigrate(&model.User{})

	return &PostgresRepo{db: db}
}


func (p *PostgresRepo) Save(user *model.User) error {
	return p.db.Create(user).Error
}

func (p *PostgresRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := p.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresRepo) Update(user *model.User) error {
	return p.db.Save(user).Error
}

func (p *PostgresRepo) Delete(user *model.User) error {
	return p.db.Delete(user).Error
}

func (p *PostgresRepo) GetAll() ([]model.User, error) {
	var users []model.User
	if err := p.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (p *PostgresRepo) RecoverPassword(email string) error {
	var user model.User
	if err := p.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	hashedPassword, _ := hash.HashPassword(email)
	return p.db.Model(&model.User{}).Where("email = ?", email).Update("password", hashedPassword).Error
}


