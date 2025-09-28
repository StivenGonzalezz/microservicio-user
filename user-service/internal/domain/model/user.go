package model

type User struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	LastName  string
	PhoneNumber string
	Email     string `gorm:"unique"`
	Password  string
	IsActive  bool
	CreatedAt string
}