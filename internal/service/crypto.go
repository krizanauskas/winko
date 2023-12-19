package service

import (
	"fmt"
	"github.com/krizanauskas/winko/pkg/bscclient"
)

type SendBnbToAddressParams struct {
	Address string `json:"recipient_address"`
	Amount  string `json:"amount"`
}

const GetBnbAllocationFunctionName = "get_bnb_allocation"
const SendBnbToAddressFunctionName = "send_bnb_to_address"

type CryptoServiceI interface {
	GetBnbAllocation(address string) (string, error)
	SendBnbToAddress(address string, amount string) (string, error)
}

type CryptoService struct {
	bscClient bscclient.ClientI
}

func (s *CryptoService) SendBnbToAddress(address string, amount string) (string, error) {
	return fmt.Sprintf("Sending %s BNB to %s via blockchain...", amount, address), nil
}

func NewCryptoService(bscClient bscclient.ClientI) *CryptoService {
	return &CryptoService{
		bscClient: bscClient,
	}
}

func (s *CryptoService) GetBnbAllocation(address string) (string, error) {
	return s.bscClient.GetBnbAllocation(address)
}
