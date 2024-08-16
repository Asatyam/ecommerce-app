package data

import (
	"database/sql"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
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

func ValidateUser(v *validator.Validator, user *User) {

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(govalidator.IsEmail(user.Email), "email", "email is not valid")

	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must be less than or equal to 500")

	v.Check(*user.Password.plaintext != "", "password", "must be provided")
	v.Check(len(*user.Password.plaintext) >= 8, "password", "must be of at least 8 characters")
	v.Check(len(*user.Password.plaintext) <= 72, "password", "must be of at most 72 characters")

	if nil == user.Password.hash {
		panic("User Password.hash is required")
	}
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
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
