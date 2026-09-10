package main

import (
	"os"

	"code.siru.ink/vdh/wahlen-web/internal/db"
	"code.siru.ink/vdh/wahlen-web/internal/routing"
)

func main() {
	db, err := db.InitializeDb("data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "80"
	}

	router := routing.InitializeRouter(db)

	router.Run(":" + port)
}
