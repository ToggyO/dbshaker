package main

import (
	main2 "github.com/ToggyO/dbshaker/cmd/dbshaker"
	"log"
)

func main() {
	if err := main2.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
