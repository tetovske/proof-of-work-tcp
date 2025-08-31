package tcp

import (
	"bufio"
	"encoding/json"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/tetovske/proof-of-work-tcp/internal/config"
	"github.com/tetovske/proof-of-work-tcp/internal/model"
	"github.com/tetovske/proof-of-work-tcp/pkg/hashcash"
)

type QuotesRepository interface {
	GetRandom() *model.Quote
}

type Transport struct {
	cfg        *config.Config
	quotesRepo QuotesRepository
}

func New(cfg *config.Config, quotesRepo QuotesRepository) *Transport {
	return &Transport{
		cfg:        cfg,
		quotesRepo: quotesRepo,
	}
}

func (t *Transport) Serve(conn net.Conn) error {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error().Msgf("failed to close tcp conn: %v", err)
		}
	}()

	connReader := bufio.NewReader(conn)

	ch, err := hashcash.CreateChallenge(t.cfg.Hashcash.Complexity, t.cfg.Hashcash.Secret)
	if err != nil {
		log.Error().Msgf("failed to create challenge: %v", err)
		return err
	}

	log.Info().Msg("challenge has been created")

	if _, err = conn.Write(ch); err != nil {
		log.Error().Msgf("failed to write data: %v", err)
		return err
	}

	log.Info().Msgf("challenge has been sent to client")

	nonce := make([]byte, hashcash.NonceSize)
	if _, err = connReader.Read(nonce); err != nil {
		log.Error().Msgf("failed to read data: %v", err)
		return err
	}

	log.Info().Msg("received nonce from client")

	if err = hashcash.Validate(t.cfg.Hashcash.Secret, ch, nonce, t.cfg.Hashcash.TTL); err != nil {
		log.Info().Msgf("challenge validation failed: %v", err)
		return nil
	}

	quote := t.quotesRepo.GetRandom()
	if quote == nil {
		log.Error().Msg("quote repo is empty")
		return nil
	}

	quoteMsg, err := json.Marshal(quote)
	if err != nil {
		log.Error().Msgf("failed to marshal quote msg: %v", err)
		return err
	}
	quoteMsg = append(quoteMsg, '\n')

	if _, err = conn.Write(quoteMsg); err != nil {
		log.Error().Msgf("failed to write data: %v", err)
		return err
	}

	log.Info().Msg("quote has been sent to client")

	return nil
}
