package config

import "time"

// RateLimitConfig is the data needed for constructing a rate limit for the HTTP requests.
type RateLimitConfig struct {
	// Burst represents the number of requests that we limit over a period frequency.
	Burst *int `yaml:"burst"`

	// Period is the number of times to allow a burst per second.
	Period *time.Duration `yaml:"period"`
}

func (rl RateLimitConfig) validate() error {
	if rl.Burst == nil {
		return MissingRateLimitFieldError("burst")
	}

	if rl.Period == nil {
		return MissingRateLimitFieldError("period")
	}

	return nil
}
