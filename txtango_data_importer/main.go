package main

import (
	"log"
	"tx2db/cmd"

	"github.com/joho/godotenv"
)

func main() {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
