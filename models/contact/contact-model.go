package contact_model

import (
	"errors"
	"paper_back/config"
)

type Contact struct {
	Id          int64  `json:"id"`
	InitiatorId int64  `json:"initiator_id"`
	AgreeingId  int64  `json:"agreeing_id"`
	Secret      string `json:"secret"`
	Status      string `json:"status"`
	CreatedAt   []byte `json:"created_at"`
	UpdatedAt   []byte `json:"updated_at"`
}

func CreateContact(initiator int64, agreeing int64) (Contact, error) {
	var userIndex int64
	user := config.Database.QueryRow(`
		SELECT id
		FROM contacts
		WHERE initiator_id = ? OR agreeing_id = ?
	`, initiator, agreeing)

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
		VALUES (?, ?, ?)
	`, initiator, agreeing, "", "ACTIVE")

	if err := result.Err(); err != nil {
		return Contact{}, errors.New("Contact was't create")
	}

	var newContactIndex int64

	newContact := config.Database.QueryRow(`
		SELECT id
		FROM contacts
		WHERE agreeing_id = ? AND initiator_id = ?
	`, agreeing, initiator)

	if err := newContact.Scan(&newContactIndex); err != nil {
		return Contact{}, errors.New("User was created")
	}

	return Contact{AgreeingId: agreeing, InitiatorId: initiator, Status: "ACTIVE", Id: newContactIndex}, nil
}

func GetContacts(userId int64, limit int, offset int) ([]Contact, error) {
	query, err := config.Database.Query(`
		SELECT 
			id,
			initiator_id,
			agreeing_id,
			secret,
			status,
			created_at,
			updated_at
		FROM contacts
		WHERE initiator_id = ? OR agreeing_id = ?
		ORDER BY updated_at DESC
		LIMIT ?
		OFFSET ?
	`, userId, userId, limit, offset)
	if err != nil {
		return []Contact{}, err
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
		); err != nil {
			return []Contact{}, err
		}

		contacts = append(contacts, contact)
	}

	if err := query.Err(); err != nil {
		return []Contact{}, err
	}

	return contacts, nil
}
