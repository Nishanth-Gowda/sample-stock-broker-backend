package repository

import "broker-backend/internal/domain/model"

// UserRepository abstracts persistence of users.
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id int64) (*model.User, error)
}
