package repository

import (
	"database/sql"

	model "forum/internal/model"
)

type UserRepository interface {
	CreateUser(*model.User) (int, error)
	GetUser(int) (*model.User, error)
	GetUserByEmail(string) (*model.User, error)
	GetAllUsers() (*[]model.User, error)         // may be don't need
	UpdateUser(*model.User) (*model.User, error) // may be don't need
	DeleteUser(int) error                        // may be don't need
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(user *model.User) (int, error) {
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

func (r *userRepo) GetUser(id int) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := model.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetAllUsers() (*[]model.User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	users := []model.User{}
	for rows.Next() {
		user := model.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (r *userRepo) DeleteUser(id int) error {
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := model.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateUser(*model.User) (*model.User, error) {
	return nil, nil
}
