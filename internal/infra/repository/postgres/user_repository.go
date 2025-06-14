package postgres

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

const userInsertQuery = `INSERT INTO users(email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id`

const userSelectByEmail = `SELECT id, email, password_hash, created_at FROM users WHERE email=$1 LIMIT 1`

const userSelectByID = `SELECT id, email, password_hash, created_at FROM users WHERE id=$1 LIMIT 1`

// Ensure implementation.
var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.QueryRow(userInsertQuery, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Get(&u, userSelectByEmail, email); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var u model.User
	if err := r.db.Get(&u, userSelectByID, id); err != nil {
		return nil, err
	}
	return &u, nil
}
