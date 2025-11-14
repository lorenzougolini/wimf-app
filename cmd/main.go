package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lorenzougolini/wimf-app/internal/server"
	"github.com/lorenzougolini/wimf-app/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// init logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.Infof("Initializing application")

	port := 3001

	// start database
	logger.Println("Creating database..")
	logger.Println("Database filename at %s", cfg.DB.Filename)
	dbconn, err := sql.Open("sqlite3", cfg.DB.Filename+"?_foreign_keys=1")
	if err != nil {
		logger.WithError(err).Error("error opening SQlite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = dbconn.Close()
	}()
	db, err := database.New(dbconn)
	if err != nil {
		logger.Error("error creating AppDatabase")
		return fmt.Errorf("creating AppDatabase: %w", err)
	}
	// guestDb := database.NewGuestStore(logger)
	// guestDb.AddGuest(database.Guest{Name: "Sigrid", Email: "sig@fake-email.no"})

	// Start API server
	logger.Info("initializing API server")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	// Create API router
	srv, err := server.NewServer(logger, port, guestDb)
	if err != nil {
		logger.Fatalf("Error when creating server: %s", err)
		os.Exit(1)
	}
	if err := srv.Start(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
	return nil
}
