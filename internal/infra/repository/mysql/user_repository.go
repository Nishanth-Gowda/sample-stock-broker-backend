package mysql

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Ensure implementation.
var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	result, err := r.db.Exec("INSERT INTO users(email, password_hash, created_at) VALUES (?, ?, ?)", user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Get(&u, "SELECT id, email, password_hash, created_at FROM users WHERE email=? LIMIT 1", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var u model.User
	if err := r.db.Get(&u, "SELECT id, email, password_hash, created_at FROM users WHERE id=? LIMIT 1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
