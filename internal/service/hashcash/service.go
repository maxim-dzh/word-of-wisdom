package hashcash

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/maxim-dzh/word-of-wisdom/internal/domain"
)

const (
	defaultVersion     = 1
	headerStringFormat = "%d:%d:%s:%s::%s:%s"
)

type service struct {
}

// CalculateCounter finds the counter which is needed for getting the hash
// that satisfies the requirements
func (s *service) CalculateCounter(ctx context.Context, header *domain.HashcashHeader) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if s.CheckBits(header.Bits, sha256.Sum256([]byte(s.FormatHeader(header)))) {
				return
			}
			if header.Counter == math.MaxUint {
				return fmt.Errorf("failed to calculate the counter")
			}
			header.Counter++
		}
	}
}

// FormatHeader returns the header as a string
func (s *service) FormatHeader(header *domain.HashcashHeader) string {
	return fmt.Sprintf(
		headerStringFormat,
		header.Version,
		header.Bits,
		strconv.FormatInt(header.Timestamp, 10),
		header.Resource,
		header.Random,
		base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(header.Counter, 10))),
	)
}

// CheckBits checks if zero bits amount in result isn't less than the bitsAmount param value
func (s *service) CheckBits(bitsAmount uint, resultHash [32]byte) bool {
	resultNum := big.NewInt(0)
	resultNum.SetBytes(resultHash[:])
	targetNum := big.NewInt(1)
	targetNum.Lsh(targetNum, 256-bitsAmount)
	return resultNum.Cmp(targetNum) == -1
}

// ParseString builds a header structure from the string value
func (s *service) ParseString(header string) (*domain.HashcashHeader, error) {
	headerParts := strings.Split(header, ":")
	if len(headerParts) != 7 {
		return nil, fmt.Errorf("invalid header")
	}
	version, err := strconv.Atoi(headerParts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}
	bits, err := strconv.Atoi(headerParts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse bits parameter: %w", err)
	}
	counterBytes, err := base64.StdEncoding.DecodeString(headerParts[6])
	if err != nil {
		return nil, fmt.Errorf("failed to decode counter parameter from base64: %w", err)
	}
	counter, err := strconv.ParseUint(string(counterBytes), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode counter parameter: %w", err)
	}
	timestamp, err := strconv.ParseInt(headerParts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date parameter: %w", err)
	}
	return &domain.HashcashHeader{
		Version:   version,
		Bits:      uint(bits),
		Timestamp: timestamp,
		Resource:  headerParts[3],
		Random:    headerParts[5],
		Counter:   counter,
	}, nil
}

// GenerateHeader returns a new hashcash header
func (s *service) GenerateHeader(bits uint, resource string) (*domain.HashcashHeader, error) {
	// get random value
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return &domain.HashcashHeader{
		Version:   defaultVersion,
		Bits:      bits,
		Timestamp: time.Now().Unix(),
		Resource:  resource,
		Random:    base64.StdEncoding.EncodeToString(bytes),
	}, nil
}

// NewService returns a new instance of the hashcash service
func NewService() *service {
	return &service{}
}
