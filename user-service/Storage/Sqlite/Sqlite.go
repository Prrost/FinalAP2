package Sqlite

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"user-service/config"
	"user-service/domain"
)

type SqliteStorage struct {
	db *sql.DB
}

var (
	ErrUserAlreadyExists = errors.New("email already in use")
)

func NewSqliteStorage(cfg *config.Config) *SqliteStorage {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL DEFAULT 'NO_PASSWORD_SET',
    isAdmin BOOLEAN NOT NULL DEFAULT FALSE
);`)

	if err != nil {
		log.Fatal(err)
	}
	return &SqliteStorage{
		db: db,
	}
}

func (s *SqliteStorage) CreateUserAdmin(user domain.User) (domain.User, error) {

	exist, err := s.IsUserExists(user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if exist {
		return domain.User{}, ErrUserAlreadyExists
	}

	stmt, err := s.db.Prepare(`INSERT INTO users (email, password, isAdmin) VALUES (?, ?, ?)`)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return domain.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}

	user.ID = uint(id)
	return user, nil
}

func (s *SqliteStorage) CreateUser(user domain.User) (domain.User, error) {

	exist, err := s.IsUserExists(user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if exist {
		return domain.User{}, ErrUserAlreadyExists
	}

	stmt, err := s.db.Prepare(`INSERT INTO users (email, isAdmin) VALUES (?, ?)`)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Email, user.IsAdmin)
	if err != nil {
		return domain.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}

	user.ID = uint(id)
	return user, nil
}

func (s *SqliteStorage) GetUserByEmail(email string) (domain.User, error) {
	stmt, err := s.db.Prepare(`SELECT * FROM users WHERE email = ?`)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	var user domain.User

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *SqliteStorage) IsUserExists(email string) (bool, error) {
	var id int
	err := s.db.QueryRow(`SELECT id FROM users WHERE email = ?`, email).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SqliteStorage) GetUserByID(id int) (domain.User, error) {
	stmt, err := s.db.Prepare(`SELECT id, email, isAdmin FROM users WHERE id = ?`)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	var user domain.User

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Email, &user.IsAdmin)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
