package server

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	log "golang.org/x/exp/slog"

	"github.com/maxim-dzh/word-of-wisdom/internal/hashcash"
)

type storage interface {
	SetChallenge(id string, header *hashcash.Header)
	GetChallenge(id string) (header *hashcash.Header, exists bool)
	DeleteChallenge(id string)
}

type wisdomWordsService interface {
	GetWord() string
}

type server struct {
	sync.RWMutex
	id                  string
	addr                string
	storage             storage
	wisdomWordsService  wisdomWordsService
	challengeComplexity uint
	challengeTimeout    time.Duration
	readTimeout         time.Duration
	logger              *log.Logger
}

// Serve contains the server logic and serves requests
func (s *server) Serve(ctx context.Context) (err error) {
	var listener net.Listener
	listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return
	}
	defer func() {
		errClose := listener.Close()
		if errClose != nil {
			s.logger.Error("failed to close the listener", "error", errClose)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var conn net.Conn
			conn, err = listener.Accept()
			if err != nil {
				s.logger.Error("failed to accept a connection", "error", err)
				continue
			}
			go s.processConn(conn)
		}
	}
}

func (s *server) processConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.logger.Error("failed to close the connection", "error", err)
		}
	}()
	reader := bufio.NewReader(conn)
	err := conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	if err != nil {
		s.logger.Warn("failed to set the read timeout", "error", err)
	}
	msg, err := reader.ReadString('\n')
	if err != nil {
		s.logger.Error("failed to read a message", "error", err)
		return
	}
	msg = strings.TrimSpace(msg)
	// if the message is empty it's a challenge request
	if msg == "" {
		// create, save and return a challenge as a hashcash header
		var challenge *hashcash.Header
		challenge, err = hashcash.NewHeader(s.challengeComplexity, s.id)
		if err != nil {
			s.logger.Error("failed to generate a hashcash header", "error", err)
			return
		}
		s.storage.SetChallenge(challenge.Random, challenge)
		err = s.writeMesssage(conn, challenge.String())
		if err != nil {
			s.logger.Error("failed to return the challenge", "error", err)
			return
		}
		return
	}
	challengeResult, err := s.parseMessage(msg)
	if err != nil {
		s.logger.Error("failed to parse the message", "error", err)
		return
	}
	originChallenge, ok := s.storage.GetChallenge(challengeResult.Random)
	if !ok {
		s.logger.Error("origin challenge not found", "error", err)
		return
	}
	defer s.storage.DeleteChallenge(challengeResult.Random)
	err = s.verify(originChallenge, challengeResult)
	if err != nil {
		s.logger.Error("the challenge failed", "error", err)
		return
	}
	// return the word
	err = s.writeMesssage(conn, s.wisdomWordsService.GetWord())
	if err != nil {
		s.logger.Error("failed to return the word of wisdom", "error", err)
		return
	}
}

func (s *server) parseMessage(msg string) (header *hashcash.Header, err error) {
	header, err = hashcash.ParseString(msg)
	if err != nil {
		err = fmt.Errorf("failed to parse string header: %w", err)
		return
	}
	return
}

// verify checks the challenge solution for correctness
// and if it's satisfying the requirements
func (s *server) verify(challengeHeader, challengeResult *hashcash.Header) (err error) {
	// we make counters equal, in order to compare other fields
	// by comparing two strings
	counter := challengeResult.Counter
	challengeResult.Counter = 0
	challengeHeader.Counter = 0
	if challengeResult.String() != challengeHeader.String() {
		return fmt.Errorf("invalid challenge result")
	}
	if time.Now().Unix()-challengeHeader.Timestamp > int64(s.challengeTimeout.Seconds()) {
		return fmt.Errorf("challenge timeout expired")
	}
	// return the resulting counter and check zero bits amount
	challengeResult.Counter = counter
	hash := sha256.Sum256([]byte(challengeResult.String()))
	if !hashcash.BitsAmountIsCorrect(challengeHeader.Bits, hash) {
		return fmt.Errorf("the challenge result isn't correct")
	}
	s.logger.Info("the challenge has resolved", "hash", fmt.Sprintf("%x", hash))
	return nil
}

func (s *server) writeMesssage(conn net.Conn, msg string) (err error) {
	_, err = fmt.Fprintf(conn, "%s\n", msg)
	return err
}

// NewServer returns a new server instance
func NewServer(
	id string,
	addr string,
	challengeComplexity uint,
	challengeTimeout time.Duration,
	readTimeout time.Duration,
	storage storage,
	wisdomWordsService wisdomWordsService,
	logger *log.Logger,
) *server {
	return &server{
		id:                  id,
		addr:                addr,
		challengeComplexity: challengeComplexity,
		challengeTimeout:    challengeTimeout,
		readTimeout:         readTimeout,
		storage:             storage,
		wisdomWordsService:  wisdomWordsService,
		logger:              logger,
	}
}
