package main

import (
	"context"
	"log"
	"os"

	app "awesomeProject/internal/api"
)

// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes http
// @BasePath /
func main() {

	log.Println("Application start!")

	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Println("can't create application")
		log.Println(err)
		os.Exit(2)
	}

	log.Println("Application start!")
	application.StartServer()
	log.Println("Application terminated!")

}
