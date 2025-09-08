package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) SetPassword(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12) // 12 is security standard balance performance and secure
	if err != nil {
		return err
	}
	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) CheckPassword(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
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

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope string, plainTextPassword string) (*User, error)
}

func (us *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, bio) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at, updated_at
	`
	err := us.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (us *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at 
	FROM users 
	WHERE username = $1
	`
	err := us.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *PostgresUserStore) UpdateUser(user *User) error {
	_, err := us.GetUserByUsername(user.Username)
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	} else if err != nil {
		return err
	}

	query := `
	UPDATE users 
	SET email = $1, bio = $2, updated_at = CURRENT_TIMESTAMP
	WHERE username = $3 
	RETURNING updated_at
	`
	result, err := us.db.Exec(query, user.Email, user.Bio, user.Username)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (us *PostgresUserStore) GetUserToken(scope string, plainTextToken string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextToken))
	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.bio, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.scope = $1 AND t.hash = $2 AND t.expiry > now()
	`
	user := &User{PasswordHash: password{}}
	err := us.db.QueryRow(query, scope, tokenHash[:]).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
