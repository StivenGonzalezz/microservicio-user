package ports

import "user-service/internal/domain/model"

type PaginatedResult struct {
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"totalPages"`
	Data       []model.User `json:"data"`
}

type UserRepository interface {
	Save(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetId(userId int) (*model.User, error)
	Update(user *model.User) error
	Delete(user *model.User) error
	RecoverPassword(email string) error
	GetByName(nameOrEmail string) ([]model.User, error)
	GetAll() ([]model.User, error)
	GetWithPagination(name string, page, limit int, sort string) (*PaginatedResult, error)
}
