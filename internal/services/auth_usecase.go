package services

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"broker-backend/pkg/auth"
	"context"
	"errors"
	"log"
	"time"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authUsecase struct {
	repo repository.UserRepository
	jwt  auth.JWTManager
	pwd  auth.PasswordHasher
}

func NewAuthUsecase(repo repository.UserRepository, jwt auth.JWTManager, pwd auth.PasswordHasher) AuthUsecase {
	return &authUsecase{repo: repo, jwt: jwt, pwd: pwd}
}

func (u *authUsecase) SignUp(ctx context.Context, email, password string) (*model.User, error) {
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
		log.Printf("auth_service: failed to create user %s: %v", email, err)
		return nil, err
	}
	log.Printf("auth_service: user %d signed up", user.ID)
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
	token, err := u.jwt.Generate(user.ID)
	if err != nil {
		log.Printf("auth_service: failed to generate token for user %d: %v", user.ID, err)
		return "", err
	}
	log.Printf("auth_service: user %d logged in", user.ID)
	return token, nil
}
