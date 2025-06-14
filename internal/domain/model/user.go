package model

import "time"

// User represents a platform user.
// In Clean Architecture, entities contain only minimal fields and no framework dependencies.
// Validation or DB specifics live in outer layers.

type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"` // hashed password stored
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
