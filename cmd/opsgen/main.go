package main

import (
	"os"

	"github.com/adnsio/gbemu/pkg/opsgen"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	opg, err := opsgen.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := opg.Generate(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
