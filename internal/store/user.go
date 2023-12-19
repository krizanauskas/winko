package store

import (
	"github.com/krizanauskas/winko/models"
)

type UserStoreI interface {
	FindUserByTgId(id int64) (*models.User, error)
}

type UserStore struct{}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (us *UserStore) FindUserByTgId(id int64) (*models.User, error) {
	threadId := "thread_dh7cXcwMJ42fQEzJyUNT6UXU"
	cryptoWalletAddress := "0xAf0476C27A15b2A6C7b9BDFe410fe0E59Ef7bEAA"

	return &models.User{
		Id:                  1,
		TgUserId:            504459620,
		LastThreadId:        &threadId,
		CryptoWalletAddress: &cryptoWalletAddress,
	}, nil
}
