package server

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/tetovske/proof-of-work-tcp/internal/config"
)

type TCPTransport interface {
	Serve(conn net.Conn) error
}

type Server struct {
	cfg *config.Config

	transport TCPTransport
}

func New(cfg *config.Config, transport TCPTransport) *Server {
	return &Server{
		cfg:       cfg,
		transport: transport,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", ":"+s.cfg.Server.Port)
	if err != nil {
		log.Error().Msgf("failed to listen tcp port: %v", err)

		return err
	}

	defer func() {
		if err = ln.Close(); err != nil {
			log.Error().Msgf("failed to close listener: %v", err)
		}
	}()

	wg := &sync.WaitGroup{}

	go func() {
		for {
			conn, aErr := ln.Accept()
			if aErr != nil {
				if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
					break
				}

				log.Error().Msgf("failed to accept: %v", aErr)
				continue
			}

			log.Info().Msg("accepting connection")

			wg.Add(1)
			go func() {
				defer wg.Done()

				s.serveConn(conn)
			}()
		}
	}()

	<-ctx.Done()

	log.Info().Msg("stopping server...")

	wg.Wait()

	return nil
}

func (s *Server) serveConn(conn net.Conn) {
	defer func() {
		r := recover()
		if r != nil {
			log.Error().Msgf("recover from panic: %v", r)
		}
	}()

	if err := s.transport.Serve(conn); err != nil {
		log.Error().Msgf("failed to serve conn: %v", err)
	}
}
