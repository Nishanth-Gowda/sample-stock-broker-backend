package auth

import "golang.org/x/crypto/bcrypt"

// PasswordHasher generates and verifies password hashes.

type PasswordHasher interface {
	Hash(pwd string) (string, error)
	Verify(hashedPwd, plainPwd string) error
}

type bcryptHasher struct{ cost int }

func NewPasswordHasher(cost int) PasswordHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

func (b *bcryptHasher) Hash(pwd string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), b.cost)
	return string(h), err
}

func (b *bcryptHasher) Verify(hashedPwd, plainPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
}
