package store

import "github.com/krizanauskas/winko/models"

type AssistantStoreI interface {
	GetAssistant() (models.Assistant, error)
}

type AssistantStore struct{}

func NewAssistantStore() *AssistantStore {
	return &AssistantStore{}
}

func (as *AssistantStore) GetAssistant() (models.Assistant, error) {
	return models.Assistant{
		Id: "asst_SlrV2qbUDOWtNESqpi2uqyOD",
	}, nil
}
