package main

import (
	"code.siru.ink/vdh/wahlen-web/internal/db"
	"code.siru.ink/vdh/wahlen-web/internal/routing"
)

func main() {
	db, err := db.InitializeDb("data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := routing.InitializeRouter(db)

	router.Run(":80")
}
