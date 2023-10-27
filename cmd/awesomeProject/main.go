package main

import (
	"log"
	"os"

	app "awesomeProject/internal/api"
)

func main() {

	log.Println("Application start!")

	application, err := app.New()
	if err != nil {
		log.Println("can't create application")
		os.Exit(2)
	}

	log.Println("Application start!")
	application.StartServer()
	log.Println("Application terminated!")

}
