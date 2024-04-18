package user_service

import (
	"errors"
	user_model "paper_back/models/user"
	token_service "paper_back/services/token"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetUser(id int64) (user_model.User, error) {
	user, err := user_model.GetUserById(id)

	if err != nil {
		return user_model.User{}, err
	}

	return user, nil
}

func Registration(username string, password string, role string) (string, error) {
	validUsername := strings.Trim(username, " ")
	validPassword := strings.Trim(password, " ")
	validRole := strings.Trim(role, " ")

	if validUsername == "" || len(validUsername) < len(username) {
		return "", errors.New("invalid username")
	}
	if validPassword == "" || len(validPassword) < len(password) {
		return "", errors.New("invalid password")
	}
	if validRole == "" || len(validRole) < len(role) {
		return "", errors.New("invalid role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validPassword), 10)
	if err != nil {
		return "", errors.New("error hash password")
	}

	user, err := user_model.CreateUser(validUsername, string(hashedPassword), validRole)
	if err != nil {
		return "", err
	}

	payload := token_service.Payload{Id: user.Id, Username: user.Username, Role: user.Role}
	tokenFromDatabase, err := token_service.CreateToken(payload)
	if err != nil {
		return "", err
	}

	return tokenFromDatabase.Token, nil
}

func Login(username string, password string) (string, error) {
	validUsername := strings.Trim(username, " ")
	validPassword := strings.Trim(password, " ")

	if validUsername == "" || len(validUsername) < len(username) {
		return "", errors.New("invalid username or password")
	}
	if validPassword == "" || len(validPassword) < len(password) {
		return "", errors.New("invalid username or password")
	}

	user, err := user_model.GetUserByUsername(validUsername)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(validPassword)); err != nil {
		return "", errors.New("invalid username or password")
	}

	payload := token_service.Payload{Id: user.Id, Username: user.Username, Role: user.Role}
	tokenFromDatabase, err := token_service.CreateToken(payload)
	if err != nil {
		return "", errors.New("error create token")
	}

	return tokenFromDatabase.Token, nil
}
