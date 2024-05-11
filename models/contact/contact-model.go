package contact_model

import (
	"errors"
	"paper_back/config"
	"paper_back/utils"
)

type Contact struct {
	Id          int64  `json:"id"`
	InitiatorId int64  `json:"initiator_id"`
	AgreeingId  int64  `json:"agreeing_id"`
	Secret      string `json:"secret"`
	Status      string `json:"status"`
	CreatedAt   []byte `json:"created_at"`
	UpdatedAt   []byte `json:"updated_at"`
	Username    string `json:"username"`
}

func CreateUserCode(userId int64) (string, error) {
	code := utils.GenerateCode()
	_, err := config.Database.Exec(`
	INSERT INTO user_codes(
		user_id,
		code
	) VALUES (?, ?)
	`, userId, code)
	if err != nil {
		config.Database.Exec(`
		UPDATE user_codes
		SET code = ?
		WHERE user_id = ?
		`, code, userId)
	}

	return code, nil
}

func CreateContact(initiator int64, code string) (Contact, error) {
	agreeingRow := config.Database.QueryRow(`
		SELECT user_id as agreeing
		FROM user_codes
		WHERE code = ?
	`, code)
	var agreeing int64

	agreeingRow.Scan(&agreeing)

	var userIndex int64
	user := config.Database.QueryRow(`
		SELECT id
		FROM contacts
		WHERE (initiator_id = ? AND agreeing_id = ?) OR (initiator_id = ? AND agreeing_id = ?)
	`, initiator, agreeing, agreeing, initiator)

	user.Scan(&userIndex)

	if userIndex > 0 {
		return Contact{}, errors.New("Contact already exists")
	}

	result := config.Database.QueryRow(`
		INSERT INTO contacts (
			initiator_id, 
			agreeing_id,
			secret,
			status
		)
		VALUES (?, ?, ?, ?)
	`, initiator, agreeing, utils.GenerateCode(), "ACTIVE")

	if err := result.Err(); err != nil {
		return Contact{}, err
	}

	var newContactIndex int64
	var username string

	newContact := config.Database.QueryRow(`
		SELECT contacts.id as id, users.username as username
		FROM contacts
		JOIN users on users.id = agreeing_id
		WHERE agreeing_id = ? AND initiator_id = ?
	`, agreeing, initiator)

	if err := newContact.Scan(&newContactIndex, &username); err != nil {
		return Contact{}, err
	}

	return Contact{AgreeingId: agreeing, InitiatorId: initiator, Status: "ACTIVE", Id: newContactIndex, Username: username}, nil
}

func GetContacts(userId int64, limit int, offset int) ([]Contact, int64, error) {
	query, err := config.Database.Query(`
		SELECT 
			contacts.id as id,
			initiator_id,
			agreeing_id,
			secret,
			status,
			contacts.created_at as created_at,
			contacts.updated_at as updated_at,
			users.username
		FROM contacts
		JOIN users on (users.id = initiator_id OR users.id = agreeing_id) AND users.id != ?
		WHERE (initiator_id = ? OR agreeing_id = ?)
		ORDER BY updated_at DESC
		LIMIT ?
		OFFSET ?
	`, userId, userId, userId, limit, offset)
	if err != nil {
		return []Contact{}, 0, err
	}

	countQuery := config.Database.QueryRow(`
		SELECT count(contacts.id) as count
		FROM contacts
		JOIN users on (users.id = initiator_id OR users.id = agreeing_id) AND users.id != ?
		WHERE (initiator_id = ? OR agreeing_id = ?)
	`, userId, userId, userId)
	if err != nil {
		return []Contact{}, 0, err
	}

	var count int64

	if err = countQuery.Scan(&count); err != nil {
		return []Contact{}, 0, err
	}

	var contacts []Contact

	for query.Next() {
		var contact Contact
		if err := query.Scan(
			&contact.Id,
			&contact.InitiatorId,
			&contact.AgreeingId,
			&contact.Secret,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.Username,
		); err != nil {
			return []Contact{}, 0, err
		}

		contacts = append(contacts, contact)
	}

	if err := query.Err(); err != nil {
		return []Contact{}, 0, err
	}

	return contacts, count, nil
}

func GetOneContact(contactId int64) (Contact, error) {
	query := config.Database.QueryRow(`
		SELECT 
			contacts.id as id,
			initiator_id,
			agreeing_id,
			secret,
			status,
			contacts.created_at as created_at,
			contacts.updated_at as updated_at,
			users.username
		FROM contacts
		JOIN users on users.id = initiator_id OR users.id = agreeing_id
		WHERE contacts.id = ?
		ORDER BY updated_at DESC
		LIMIT 1
	`, contactId)

	var contact Contact

	if err := query.Scan(
		&contact.Id,
		&contact.InitiatorId,
		&contact.AgreeingId,
		&contact.Secret,
		&contact.Status,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.Username,
	); err != nil {
		return Contact{}, err
	}

	return contact, nil
}
