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
)

const (
	defaultVersion     = 1
	headerStringFormat = "%d:%d:%s:%s::%s:%s"
)

// Header is an object representation of the hashcash header
type Header struct {
	Version   int
	Bits      uint
	Timestamp int64
	Resource  string
	Random    string
	Counter   uint64
}

// CalculateCounter finds the counter which is needed for getting the hash
// that satisfies the requirements
func (h *Header) CalculateCounter(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if BitsAmountIsCorrect(h.Bits, sha256.Sum256([]byte(h.String()))) {
				return
			}
			if h.Counter == math.MaxUint {
				return fmt.Errorf("failed to calculate the counter")
			}
			h.Counter++
		}
	}
}

// String returns the header as a string
func (h *Header) String() string {
	return fmt.Sprintf(
		headerStringFormat,
		h.Version,
		h.Bits,
		strconv.FormatInt(h.Timestamp, 10),
		h.Resource,
		h.Random,
		base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(h.Counter, 10))),
	)
}

// BitsAmountIsCorrect checks if zero bits amount in result isn't less than the bitsAmount param
func BitsAmountIsCorrect(bitsAmount uint, resultHash [32]byte) bool {
	resultNum := big.NewInt(0)
	resultNum.SetBytes(resultHash[:])
	targetNum := big.NewInt(1)
	targetNum.Lsh(targetNum, 256-bitsAmount)
	return resultNum.Cmp(targetNum) == -1
}

// ParseString build a header structure from the string value
func ParseString(header string) (*Header, error) {
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
	return &Header{
		Version:   version,
		Bits:      uint(bits),
		Timestamp: timestamp,
		Resource:  headerParts[3],
		Random:    headerParts[5],
		Counter:   counter,
	}, nil
}

// NewHeader returns a new hashcash header
func NewHeader(bits uint, resource string) (*Header, error) {
	// get random value
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return &Header{
		Version:   defaultVersion,
		Bits:      bits,
		Timestamp: time.Now().Unix(),
		Resource:  resource,
		Random:    base64.StdEncoding.EncodeToString(bytes),
	}, nil
}
