package models

type User struct {
	Id                  int64
	TgUserId            int64
	LastThreadId        *string
	CryptoWalletAddress *string
}
