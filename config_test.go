// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Parallel()
	t.Run("Both empty", func(t *testing.T) {
		t.Parallel()
		rlc := RateLimitConfig{}
		result := rlc.validate()
		if errors.Is(result, MissingRateLimitFieldError("burst")) {
			t.Errorf("expected: %v, got: %v", MissingRateLimitFieldError("burst"), rlc.validate())
		}
	})
	t.Run("Period Empty", func(t *testing.T) {
		t.Parallel()
		rlc := RateLimitConfig{
			Burst:  1,
			Period: 0,
		}
		if errors.Is(rlc.validate(), MissingRateLimitFieldError("period")) {
			t.Errorf("expect: %v, got: %v", MissingRateLimitFieldError("period"), rlc.validate())
		}
	})
	t.Run("Burst Empty", func(t *testing.T) {
		t.Parallel()
		rlc := RateLimitConfig{
			Burst:  0,
			Period: 1,
		}
		if errors.Is(rlc.validate(), MissingRateLimitFieldError("burst")) {
			t.Errorf("expect: %v, got: %v", MissingRateLimitFieldError("burst"), rlc.validate())
		}
	})
	t.Run("Valid input", func(t *testing.T) {
		t.Parallel()
		rlc := RateLimitConfig{
			Burst:  1,
			Period: 1,
		}
		if rlc.validate() != nil {
			t.Errorf("expect: %v, got: %v", nil, rlc.validate())
		}
	})
}
