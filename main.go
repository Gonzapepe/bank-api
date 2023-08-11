package main

import (
	"log"

	"github.com/Gonzapepe/bank-api/api"
	"github.com/Gonzapepe/bank-api/storage"
)

func main() {
	store, err := storage.NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", store)
	server.Run()
}