package config

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Configuration struct {
	FiberConfig fiber.Config
}

// Verify Startup Enviromental Vars and setup the Fiber option
func Startup() Configuration {
	var conf Configuration

	requiredEnvVars := []string{
		"API_KEY",
		"NDAYS",
		"SYMBOL",
	}

	for index := range requiredEnvVars {
		if os.Getenv(requiredEnvVars[index]) == "" {
			log.Fatalf("Missing %s Env var \n", requiredEnvVars[index])
		}
	}

	// Fiber Options Set statically for now.
	conf.FiberConfig = fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Microservice-Test",
		AppName:       "MicroService Test v0.01",
		ReadTimeout:   (30 * time.Second),
	}

	return conf
}
