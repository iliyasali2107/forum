package repository

import (
	"database/sql"

	models "forum/internal/models"
)

type UserRepository interface {
	CreateUser(*models.User) (int, error)
	GetUser(int) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByToken(token string) (*models.User, error)
	GetAllUsers() (*[]models.User, error)          // may be don't need
	UpdateUser(*models.User) (*models.User, error) // may be don't need
	DeleteUser(int) error                          // may be don't need
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) CreateUser(user *models.User) (int, error) {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	row, err := ur.db.Exec(query, user.Name, user.Email, user.Password.Hash)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (ur *userRepo) GetUser(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	row := ur.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := models.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepo) GetAllUsers() (*[]models.User, error) {
	query := `SELECT * FROM users`
	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (ur *userRepo) DeleteUser(id int) error {
	return nil
}

func (ur *userRepo) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = ?`
	row := ur.db.QueryRow(query, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := models.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepo) UpdateUser(*models.User) (*models.User, error) {
	return nil, nil
}

func (ur *userRepo) GetUserByToken(token string) (*models.User, error) {
	return nil, nil
}
