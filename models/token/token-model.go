package token_model

import (
	"errors"
	"paper_back/config"
)

type Token struct {
	Id int64 `json:"id"`
	UserId int64 `json:"user_id"`
	Token string `json:"access"`
}

func AddToken(token string, userId int64) (Token, error) {
	var userIdFromDatabase int64
	oldTokenData := config.Database.QueryRow(`
		SELECT id
		FROM tokens
		WHERE user_id = ?
	`, userId)

	if err := oldTokenData.Scan(&userIdFromDatabase); err == nil && userIdFromDatabase > 0 {
		if _, newErr := config.Database.Exec(`
			DELETE FROM tokens
			WHERE user_id = ?
		`, userId); newErr != nil {
			return Token{}, errors.New("error create token")
		}
	}

	if _, err := config.Database.Exec(`
		INSERT INTO tokens (
			user_id, 
			access
		) VALUES (
			?, ?
		)
	`, userId, token); err != nil {
		return Token{}, errors.New("error create token")
	}

	return Token{Token: token, UserId: userId}, nil
}