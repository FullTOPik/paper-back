package message_model

import (
	"errors"
	"paper_back/config"
	contact_model "paper_back/models/contact"
)

type Message struct {
	Id        int64  `json:"id"`
	ContactId int64  `json:"contact_id"`
	Sender    string `json:"sender"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func CreateMessage(userId int64, contactId int64, text string) (string, error) {
	contact, err := contact_model.GetOneContact(contactId)
	if err != nil {
		return "", err
	}

	if contact.AgreeingId != userId && contact.InitiatorId != userId {
		return "", errors.New("contact not found")
	}

	if _, err := config.Database.Exec(`
		INSERT INTO messages (
			contact_id,
			sender,
			text
		) VALUES (?, "?", ?)
	`, contactId, userId, text); err != nil {
		return "", err
	}

	return text, nil
}

func GetMessages(userId int64, contactId int64) ([]Message, error) {
	contact, err := contact_model.GetOneContact(contactId)
	if err != nil {
		return []Message{}, err
	}

	if contact.AgreeingId != userId && contact.InitiatorId != userId {
		return []Message{}, errors.New("contact not found")
	}

	var messages []Message

	dataRows, err := config.Database.Query(`
		SELECT 
			id,
			contact_id,
			sender,
			text,
			created_at,
			updated_at
		FROM messages
		WHERE contact_id = ?
	`, contactId)
	if err != nil {
		return []Message{}, err
	}

	for dataRows.Next() {
		var message Message
		if err := dataRows.Scan(
			&message.Id,
			&message.ContactId,
			&message.Sender,
			&message.Text,
			&message.CreatedAt,
			&message.UpdatedAt,
		); err != nil {
			return []Message{}, err
		}

		messages = append(messages, message)
	}

	if err := dataRows.Err(); err != nil {
		return []Message{}, err
	}

	return messages, nil
}
