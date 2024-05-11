package contact_service

import (
	"math"
	contact_model "paper_back/models/contact"
)

func GetContacts(userId int64, pageSize int, page int) ([]contact_model.Contact, int64, error) {
	offset := (page - 1) * pageSize

	contacts, count, err := contact_model.GetContacts(userId, pageSize, offset)
	if err != nil {
		return []contact_model.Contact{}, 0, err
	}

	return contacts, int64(math.Ceil(float64(count) / float64(pageSize))), nil
}

func CreateContact(currentUserId int64, code string) (contact_model.Contact, error) {
	contact, err := contact_model.CreateContact(currentUserId, code)
	if err != nil {
		return contact_model.Contact{}, err
	}

	return contact, nil
}
