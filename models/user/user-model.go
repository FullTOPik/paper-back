package user_model

import (
	"errors"
	"paper_back/config"
)

type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	LastVisit string `json:"last_visit"`
	Password  string `json:"password"`
}

func GetUserById(id int64) (User, error) {
	result := config.Database.QueryRow(`
		SELECT 
			id, 
			username, 
			role, 
			created_at
		FROM users
		WHERE id = ?
	`, id)

	var user User

	if err := result.Scan(
		&user.Id,
		&user.Username,
		&user.Role,
		&user.CreatedAt,
	); err != nil {
		return User{}, errors.New("User not found")
	}

	return user, nil
}

func CreateUser(username string, password string, role string) (User, error) {
	var userIndex int64
	user := config.Database.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?
	`, username)

	user.Scan(&userIndex)

	if userIndex > 0 {
		return User{}, errors.New("User already exists")
	}

	result := config.Database.QueryRow(`
		INSERT INTO users (
			username, 
			role,
			password
		)
		VALUES (?, ?, ?)
	`, username, role, password)

	if err := result.Err(); err != nil {
		return User{}, errors.New("User was't create")
	}

	var newUserIndex int64

	newUser := config.Database.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?
	`, username)

	if err := newUser.Scan(&newUserIndex); err != nil {
		return User{}, errors.New("User was created")
	}

	return User{Username: username, Role: role, Id: newUserIndex}, nil
}

func GetUserByUsername(username string) (User, error) {
	var user User

	userData := config.Database.QueryRow(`
		SELECT 
			id, password, role, username, created_at
		FROM users
		WHERE username = ?
	`, username)

	if err := userData.Scan(&user.Id, &user.Password, &user.Role, &user.Username, &user.CreatedAt); err != nil {
		return User{}, errors.New("invalid username")
	}

	return user, nil
}
