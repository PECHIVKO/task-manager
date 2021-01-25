package main

import (
	"log"

	"github.com/PECHIVKO/task-manager/server"
	_ "github.com/lib/pq"
)

func main() {
	app := server.NewApp()

	if err := app.Run("8181"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
