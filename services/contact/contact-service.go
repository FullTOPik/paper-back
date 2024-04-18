package contact_service

import (
	"errors"
	contact_model "paper_back/models/contact"
)

func GetContacts(userId int64, pageSize int, page int) ([]contact_model.Contact, error) {
	offset := (page - 1) * pageSize

	contacts, err := contact_model.GetContacts(userId, pageSize, offset)
	if err != nil {
		return []contact_model.Contact{}, errors.New("error to get contacts")
	}

	return contacts, nil
}