package main

import (
	"log"
	"path"
	"tx2db/cmd"

	"github.com/joho/godotenv"
	"github.com/kardianos/osext"
)

func main() {
	//get program path
	wd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	//load .env file
	err = godotenv.Load(path.Join(wd, ".env"))
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
