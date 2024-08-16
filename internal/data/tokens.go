package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type TokenModel struct {
	DB *sql.DB
}
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Scope     string    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func generateToken(userID int64, scope string, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil

}
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func (m *TokenModel) New(userID int64, scope string, ttl time.Duration) (*Token, error) {
	token, err := generateToken(userID, scope, ttl)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, err
}
func (m *TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, scope, expiry)
		VALUES ($1, $2, $3, $4)
		`
	args := []any{token.Hash, token.UserID, token.Scope, token.Expiry}
	_, err := m.DB.Exec(query, args...)
	return err
}
func (m *TokenModel) DeleteForUser(token *Token) error {

	query := `
		DELETE FROM tokens
       WHERE user_id = $1 and scope = $2 
	`
	_, err := m.DB.Exec(query, token.UserID, token.Scope)
	return err
}
