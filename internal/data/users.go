package data

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Prime     bool      `json:"prime"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
	Version   int64     `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintext
	p.hash = hash
	return nil
}

func (m UserModel) Insert(user *User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("HelloWorld"), 12)
	if err != nil {
		return err
	}
	user.Password.hash = hash
	query := `
		INSERT INTO users(name, email, password_hash, activated, prime, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version
	`
	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, user.Prime, user.IsAdmin}

	err = m.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		return err
	}
	return nil
}
