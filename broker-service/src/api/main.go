package main

import "log"

const PORT = ":8080"

type Config struct{}

func main() {
	app := Config{}

	server := app.Routes()

	err := server.Run(PORT)
	if err != nil {
		log.Panic(err)
	}
}
