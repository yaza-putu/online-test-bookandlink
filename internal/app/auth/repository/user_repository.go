package repository

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/auth/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/database"
)

// UserInterface / **************************************************************
type User interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}

type userRepository struct {
	entity entity.User
}

func NewUser() User {
	return &userRepository{
		entity: entity.User{},
	}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	e := u.entity
	r := database.Instance.Where("email = ?", email).First(&e)
	return e, r.Error
}
