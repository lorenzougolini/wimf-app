package main

import (
	"log"
	"os"

	"github.com/lorenzougolini/wimf-app/internal/server"
	"github.com/lorenzougolini/wimf-app/internal/store"
)

func main() {
	logger := log.New(os.Stdout, "[Spooktober] ", log.LstdFlags)

	port := 3001

	logger.Print("Creating guests store..")
	guestDb := store.NewGuestStore(logger)
	guestDb.AddGuest(store.Guest{Name: "Sigrid", Email: "sig@fake-email.no"})

	srv, err := server.NewServer(logger, port, guestDb)
	if err != nil {
		logger.Fatalf("Error when creating server: %s", err)
		os.Exit(1)
	}
	if err := srv.Start(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
}
