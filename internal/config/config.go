package config

import "time"

// Config is a common config for the service
type Config struct {
	ServerAddr          string        `envconfig:"SERVER_ADDR" default:":8090"`
	ChallengeTimeout    time.Duration `envconfig:"CHALLENGE_TIMEOUT" default:"5s"`
	ChallengeComplexity uint          `envconfig:"CHALLENGE_COMPLEXITY" default:"10"`
	ReadTimeout         time.Duration `envconfig:"READ_TIMEOUT" default:"30s"`
}
