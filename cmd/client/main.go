package main

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/tetovske/proof-of-work-tcp/internal/model"
	"github.com/tetovske/proof-of-work-tcp/pkg/hashcash"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	conn, err := net.Dial("tcp", os.Getenv("QUOTES_SERVER_ADDR"))
	if err != nil {
		log.Error().Msgf("failed to connect: %v", err)
		return
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Error().Msgf("failed to close connection: %v", err)
		}
	}()

	connReader := bufio.NewReader(conn)

	challenge := make([]byte, hashcash.ChallengeSize)
	if _, err = connReader.Read(challenge); err != nil {
		log.Error().Msgf("failed to read conn: %v", err)
		return
	}

	log.Info().Msg("got challenge from server")

	nonce, err := hashcash.Solve(ctx, challenge)
	if err != nil {
		log.Error().Msgf("failed to solve pow challenge: %v", err)
		return
	}

	log.Info().Msg("nonce for challenge has been found")

	if _, err = conn.Write(nonce); err != nil {
		log.Error().Msgf("failed to write to conn: %v", err)
		return
	}

	quote, err := connReader.ReadBytes('\n')
	if err != nil {
		log.Error().Msgf("failed to read conn: %v", err)
		return
	}

	var d model.Quote
	if err = json.Unmarshal(quote, &d); err != nil {
		log.Error().Msgf("failed to unmarshal msg: %v", err)
		return
	}

	log.Info().Msgf("got quote: %s", d.Text)
}
