package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/sirupsen/logrus"
)

type WebApiConfiguration struct {
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3001"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	DB struct {
		Filename string `conf:"default: ./wimf-app.db"`
	}
}

func loadConfiguration() (WebApiConfiguration, error) {
	var cfg WebApiConfiguration

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	if err := conf.Parse(os.Args[1:], "CFG", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, err := conf.Usage("CFG", &cfg)
			if err != nil {
				return cfg, fmt.Errorf("generating config usage: %w", err)
			}
			logger.Println(usage) //nolint:forbidigo
			return cfg, conf.ErrHelpWanted
		}
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
