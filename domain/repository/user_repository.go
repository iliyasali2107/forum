package repository

import (
	"context"
	"database/sql"
	"time"

	"forum/domain/models"
)

type UserRepository interface {
	CreateUser(*models.User) (int, error)
	GetUser(int) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByToken(token string) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	GetAllUsers() ([]models.User, error)           // may be no need
	UpdateUser(*models.User) (*models.User, error) // may be no need
	DeleteUser(int) error                          // may be no need

	SaveToken(*models.User) error   // from Auth Repository
	DeleteToken(token string) error // from Auth Repository
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

// TODO: UNIQUE constraint failed: users.name
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

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.Token, &user.Expires); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) GetAllUsers() ([]models.User, error) {
	query := `SELECT * FROM users`
	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *userRepo) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = ?`
	row := ur.db.QueryRow(query, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.Token, &user.Expires); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) GetUserByName(name string) (*models.User, error) {
	query := `SELECT * FROM users WHERE name = ?`
	row := ur.db.QueryRow(query, name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.Token, &user.Expires); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) GetUserByToken(token string) (*models.User, error) {
	query := `SELECT * FROM users WHERE token = ?`
	row := ur.db.QueryRow(query, token)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := &models.User{}

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.Token, &user.Expires); err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) SaveToken(user *models.User) error {
	query := `UPDATE users 
		SET token = ?,
		    expiry = ?
		WHERE id = ?;`

	args := []interface{}{user.Token, user.Expires, user.ID}
	// TODO: change *ctx* to be an argument to *repository* function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := ur.db.ExecContext(ctx, query, args...)
	return err
}

func (ur *userRepo) DeleteToken(token string) error {
	query := `UPDATE users
		SET token = NULL,
		    expiry = NULL
		WHERE token = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := ur.db.ExecContext(ctx, query, token)
	return err
}

func (ur *userRepo) UpdateUser(*models.User) (*models.User, error) {
	return nil, nil
}

func (ur *userRepo) DeleteUser(id int) error {
	return nil
}
