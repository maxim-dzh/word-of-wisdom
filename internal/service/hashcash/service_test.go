package hashcash

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/maxim-dzh/word-of-wisdom/internal/domain"
)

type hashcashTestSuite struct {
	suite.Suite

	ctx               context.Context
	err               error
	currentGoroutines goleak.Option

	// testing data
	headerString string
	header       *domain.HashcashHeader
	service      *service
}

func (s *hashcashTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.err = errors.New("error")
	s.currentGoroutines = goleak.IgnoreCurrent()
	s.headerString = "1:10:12343:b212c6f7-6a05-4bd1-b181-467e12f0cb30::WCNLFUvub2oCUA==:MA=="
	s.service = NewService()
	s.header = &domain.HashcashHeader{
		Version:   1,
		Bits:      10,
		Timestamp: 12343,
		Resource:  "b212c6f7-6a05-4bd1-b181-467e12f0cb30",
		Random:    "WCNLFUvub2oCUA==",
	}
}

func (s *hashcashTestSuite) TearDownTest() {
	goleak.VerifyNone(s.T(), s.currentGoroutines)
}

func (s *hashcashTestSuite) TestFormatHeader_Ok() {
	headerString := s.service.FormatHeader(s.header)
	s.EqualValues(s.headerString, headerString)
}

func (s *hashcashTestSuite) TestParseString_Ok() {
	header, err := s.service.ParseString(s.headerString)
	s.Nil(err)
	s.Equal(s.header, header)
}

func TestHashcashTestSuite(t *testing.T) {
	suite.Run(t, &hashcashTestSuite{})
}
