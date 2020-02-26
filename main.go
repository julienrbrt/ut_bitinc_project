package main

import (
	"log"
	"tx2db/cmd"
	"tx2db/database"

	"github.com/joho/godotenv"
)

func main() {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//connect to database
	err = database.InitDB()
	if err != nil {
		panic(err)
	}
	defer database.DB().Close()

	cmd.Execute()
}
