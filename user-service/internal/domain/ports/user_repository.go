package ports

import "user-service/internal/domain/model"

type UserRepository interface {
	Save(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetId(userId int) (*model.User, error)
	Update(user *model.User) error
	Delete(user *model.User) error
	RecoverPassword(email string) error
	GetAll(nameOrEmail string) ([]model.User, error)
}
