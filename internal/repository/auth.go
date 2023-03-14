package repository

import (
	"database/sql"

	"forum/internal/model"
)

type AuthRepository interface {
	CreateUser(*model.User) (int, error)
	GetUser(string) (*model.User, error)
	SaveToken(*model.User) error
	GetUserByToken(token string) (*model.User, error)
	DeleteToken(token string) error
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) CreateUser(user *model.User) (int, error) {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	row, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *authRepo) GetUser(name string) (*model.User, error) {
	query := `SELECT * FROM users WHERE name = ?`
	row := r.db.QueryRow(query, name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := model.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) SaveToken(user *model.User) error {
	return nil
}

func (r *authRepo) GetUserByToken(token string) (*model.User, error) {
	return nil, nil
}

func (r *authRepo) DeleteToken(token string) error {
	return nil
}
