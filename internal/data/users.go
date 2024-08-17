package data

import (
	"crypto/sha256"
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

var GuestUser = &User{}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

func (user *User) IsGuestUser() bool {
	return user == GuestUser
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

func (p *password) Matches(plaintext string) (bool, error) {

	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (m UserModel) Insert(user *User) error {

	query := `
		INSERT INTO users(name, email, password_hash, activated, prime, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version
	`
	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, user.Prime, user.IsAdmin}

	err := m.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
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

func (m UserModel) GetByEmail(email string) (*User, error) {

	query := `
		SELECT id, created_at, version, name, email, password_hash, activated, prime, is_admin
		FROM users
		WHERE email = $1
`
	var user User
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Prime,
		&user.IsAdmin,
	)
	if err != nil {
		switch {

		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil

}
func (m UserModel) GetForToken(scope string, token string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(token))
	query := `
	SELECT users.id, users.name, users.created_at,  users.email, users.activated, users.prime, users.is_admin, users.version
	FROM users
	INNER JOIN tokens ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	and tokens.scope = $2
	and tokens.expiry > $3
    `
	var user User
	err := m.DB.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.Email,
		&user.Activated,
		&user.Prime,
		&user.IsAdmin,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users 
		SET name = $1,  activated = $2, prime = $3, is_admin = $4, version = version + 1, password_hash = $5
		WHERE id=$6 AND version = $7
		RETURNING version
`
	args := []any{
		user.Name,
		user.Activated,
		user.Prime,
		user.IsAdmin,
		user.Password.hash,
		user.ID,
		user.Version,
	}
	err := m.DB.QueryRow(query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
