package store

import (
	"database/sql"
	"time"
)

type password struct {
	plainText *string
	hash      []byte
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"_"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgreUserStore struct {
	db *sql.DB
}

func NewPostgreUserStore(db *sql.DB) *PostgreUserStore {
	return &PostgreUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

func (us *PostgreUserStore) CreateUser(user *User) error {
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

func (us *PostgreUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at 
	FROM users 
	WHERE username = $1
	`
	err := us.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *PostgreUserStore) UpdateUser(user *User) error {
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
