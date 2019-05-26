package mysql

import (
	"database/sql"
	"go-api-base/model"
)

func NewUserRepository(db *sql.DB) model.UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *sql.DB
}

func (userRepository *userRepository) IsEmailExists(email string) bool {
	var exists bool

	query := `
		SELECT CASE WHEN EXISTS (
			SELECT *
			FROM user
			WHERE email = ?
		)
		THEN 1
		ELSE 0 END
	`

	userRepository.db.QueryRow(query, email).Scan(&exists)

	return exists
}

func (userRepository *userRepository) Store(user model.User) error {
	query := `
		INSERT INTO user (id, email, password, password_salt, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	_, err := userRepository.db.Exec(
		query,
		user.UserID,
		user.Email,
		user.Password.Encrypt(),
		user.Password.PasswordSalt,
		user.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (userRepository *userRepository) FindByEmail(
	email string,
) *model.User {
	var user model.User

	query := `
		SELECT user.id, user.email, user.password, user.password_salt, user.name
		FROM user
		WHERE email = ?
	`

	err := userRepository.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordEncrypted,
		&user.Password.PasswordSalt,
		&user.Name,
	)
	if err != nil {
		return nil
	}

	return &user
}
