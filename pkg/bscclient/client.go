package bscclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseUrl    string
	HttpClient *http.Client
	ApiKey     string
}

type Config struct {
	BaseURL string
	Timeout int64
	ApiKey  string
}

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type ClientI interface {
	GetBnbAllocation(address string) (string, error)
}

func NewClient(cfg Config) *Client {
	return &Client{
		BaseUrl: cfg.BaseURL,
		HttpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func (c *Client) GetBnbAllocation(address string) (string, error) {
	url := fmt.Sprintf("%s?module=account&action=balance&address=%s&apikey=%s", c.BaseUrl, address, c.ApiKey)

	resp, err := c.HttpClient.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	var apiResponse ApiResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return "", err
	}

	if apiResponse.Status != "1" {
		return "", fmt.Errorf("API response status is not 1")
	}

	allocation := apiResponse.Result.(string)

	if len(allocation) < 18 {
		allocation = strings.Repeat("0", 18-len(allocation)) + allocation
	}

	resultStr := allocation[:len(allocation)-18] + "." + allocation[len(allocation)-18:]

	if strings.HasPrefix(resultStr, ".") {
		resultStr = "0" + resultStr
	}

	return resultStr, nil
}
