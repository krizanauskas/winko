package main

import (
	"github.com/joho/godotenv"
	"github.com/krizanauskas/winko/pkg/config"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}
