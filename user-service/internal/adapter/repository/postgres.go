package repository

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"user-service/internal/domain/model"
	"user-service/internal/domain/ports"

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
	result := p.db.Create(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return errors.New("email already registered")
		}
		return result.Error
	}
	return nil
}

func (p *PostgresRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := p.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresRepo) GetId(userId int) (*model.User, error) {
	var user model.User
	if err := p.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresRepo) Update(user *model.User) error {
	return s.db.Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(user).Error
}

func (p *PostgresRepo) Delete(user *model.User) error {
	return p.db.Delete(user).Error
}

func (p *PostgresRepo) GetByName(nameOrEmail string) ([]model.User, error) {
	var users []model.User
	if err := p.db.Where("name = ? OR email = ?", nameOrEmail, nameOrEmail).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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
	return nil
}

func (p *PostgresRepo) GetWithPagination(name string, page, limit int, sort string) (*ports.PaginatedResult, error) {
	var users []model.User
	var total int64

	// Construir la consulta base
	query := p.db.Model(&model.User{})

	// Aplicar filtro por nombre si se proporciona
	if name != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?",
			"%"+name+"%",
			"%"+name+"%")
	}

	// Contar el total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calcular el offset
	offset := (page - 1) * limit

	// Aplicar ordenamiento
	orderBy := "created_at"
	if sort == "desc" {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	// Aplicar paginación y obtener los resultados
	if err := query.Order(orderBy).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	// Calcular el número total de páginas
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	// Construir y devolver el resultado paginado
	result := &ports.PaginatedResult{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       users,
	}

	return result, nil
}
