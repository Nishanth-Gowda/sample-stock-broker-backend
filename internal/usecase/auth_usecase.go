package usecase

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"broker-backend/pkg/auth"
	"context"
	"errors"
	"time"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error) // returns JWT token
}

// authUsecase implements AuthUsecase following Clean Architecture.
type authUsecase struct {
	repo repository.UserRepository
	jwt  auth.JWTManager
	pwd  auth.PasswordHasher
}

func NewAuthUsecase(repo repository.UserRepository, jwt auth.JWTManager, pwd auth.PasswordHasher) AuthUsecase {
	return &authUsecase{repo: repo, jwt: jwt, pwd: pwd}
}

func (u *authUsecase) SignUp(ctx context.Context, email, password string) (*model.User, error) {
	// Check if user exists.
	if existing, _ := u.repo.FindByEmail(email); existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := u.pwd.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     email,
		Password:  hash,
		CreatedAt: time.Now().UTC(),
	}
	if err := u.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := u.pwd.Verify(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}
	return u.jwt.Generate(user.ID)
}
