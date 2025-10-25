package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	Init()
	RunServer()
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
