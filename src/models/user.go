package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(db *sql.DB, user *User) error {
	query := `INSERT INTO users(username, password, created_at) VALUES($1, $2, $3)`
	currentTime := time.Now()

	_, err := db.Exec(query, user.Username, user.Password, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	// Определите SQL-запрос для получения данных пользователя
	query := `SELECT id, username, password, created_at FROM users WHERE username = $1`

	// Используйте QueryRow для выполнения запроса, так как ожидается одна запись
	var user User
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Обработка случая, когда пользователь не найден
			return nil, nil // Или возвращайте ошибку, если предпочитаете явно обрабатывать это как ошибочную ситуацию
		}
		// Возврат ошибки при выполнении запроса или во время сканирования
		return nil, err
	}

	return &user, nil
}
