package main

import (
	"flag"
	"log"

	"github.com/neemiasjnr/echod/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	// setup logger
	lvl := flag.String("level", "info", "Log level")
	flag.Parse()
	logLevel, err := zerolog.ParseLevel(*lvl)
	if err != nil {
		log.Fatal("Invalid log level")
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)

	app := server.New()
	log.Fatal(app.Listen(":3000"))
}
