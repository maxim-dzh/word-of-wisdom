package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/exp/slog"

	"github.com/maxim-dzh/word-of-wisdom/internal/config"
	"github.com/maxim-dzh/word-of-wisdom/internal/hashcash"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var (
		cfg    config.Config
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	)
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return
	}
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", cfg.ServerAddr)
	if err != nil {
		logger.Error("failed to connect to the server", "error", err)
		return
	}
	defer conn.Close()

	// init a challenge
	_, err = fmt.Fprintf(conn, "\n")
	if err != nil {
		logger.Error("failed to init the challenge", "error", err)
		return
	}
	// parse the challenge data
	reader := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	if err != nil {
		logger.Warn("failed to set the read timeout", "error", err)
	}
	challengeData, err := reader.ReadString('\n')
	if err != nil {
		logger.Error("failed to read the challenge data", "error", err)
		return
	}
	challenge, err := hashcash.ParseString(challengeData)
	if err != nil {
		logger.Error("failed to parse the challenge", "error", err)
		return
	}

	// calculate the counter and send the resulting header
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, cfg.ChallengeTimeout)
	defer timeoutCancel()
	err = challenge.CalculateCounter(timeoutCtx)
	if err != nil {
		logger.Error("failed to calculate counter", "error", err)
		return
	}

	// send a solution
	conn, err = dialer.DialContext(ctx, "tcp", cfg.ServerAddr)
	if err != nil {
		logger.Error("failed to connect to the server", "error", err)
		return
	}
	defer conn.Close()
	reader = bufio.NewReader(conn)
	_, err = fmt.Fprintf(conn, "%s\n", challenge.String())
	if err != nil {
		logger.Error("failed to send a challenge result", "error", err)
		return
	}

	// read the response
	err = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	if err != nil {
		logger.Warn("failed to set the read timeout", "error", err)
	}
	var result string
	result, err = reader.ReadString('\n')
	if err != nil {
		logger.Error("failed to read the word of wisdom", "error", err)
		return
	}
	logger.Info("we have received a word of wisdom", "word", result)
}
