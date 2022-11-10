package main

import (
	"log"
	"os"

	"github.com/ToggyO/dbshaker/pkg"
)

func main() {
	err := os.Mkdir("./test", os.FileMode(0777))
	if err != nil {
		log.Fatal(err)
	}

	err = dbshaker.CreateMigrationTemplate("createUsersTable", "./test", dbshaker.GoTemplate)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove("./test")
}
