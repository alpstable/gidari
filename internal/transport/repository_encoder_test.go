package transport

import (
	"errors"
	"net/url"
	"testing"
)

func TestRepositoryEncoderRegistry(t *testing.T) {
	t.Parallel()

	err := RegisterEncoders(RepositoryEncoderRegistry{}, RegisterDefaultEncoder, RegisterCBPEncoder)
	if err != nil {
		t.Fatalf("error registering encoders: %v", err)
	}

	t.Run("RepositoryEncoderKey", func(t *testing.T) {
		t.Parallel()
		t.Run("NewRepositoryEncoderKey", func(t *testing.T) {
			t.Parallel()
			testURL, _ := url.Parse("https://api.pro.coinbase.com")
			key := NewRepositoryEncoderKey(testURL)
			if key.host != "api.pro.coinbase.com" {
				t.Errorf("expected key.host to be %s, got %s", "api.pro.coinbase.com", key.host)
			}
		})
		t.Run("NewDefaultRepositoryEncoderKey", func(t *testing.T) {
			t.Parallel()
			key := NewDefaultRepositoryEncoderKey()
			if key.host != "" {
				t.Errorf("expected key.host to be %s, got %s", "", key.host)
			}
		})
	})
	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		t.Run("should return ErrRepositoryEncoderExists when registering an encoder that already exists",
			func(t *testing.T) {
				t.Parallel()
				registry := RepositoryEncoderRegistry{}
				key := RepositoryEncoderKey{host: "test"}
				registry[key] = new(DefaultRepositoryEncoder)

				testURL, _ := url.Parse("http://test")
				err := registry.Register(testURL, new(DefaultRepositoryEncoder))
				if !errors.Is(err, ErrRepositoryEncoderExists) {
					t.Errorf("expected error to be %q, got %q", ErrRepositoryEncoderExists, err)
				}
			})
		t.Run("should register an encoder when one does not already exist", func(t *testing.T) {
			t.Parallel()
			registry := RepositoryEncoderRegistry{}
			testURL, _ := url.Parse("http://test")
			err := registry.Register(testURL, new(DefaultRepositoryEncoder))
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}

			key := NewRepositoryEncoderKey(testURL)
			if registry[key] == nil {
				t.Error("expected encoder to be registered, got nil")
			}
		})
	})
	t.Run("Lookup", func(t *testing.T) {
		t.Run("should return the encoder when one is registered", func(t *testing.T) {
			t.Parallel()
			registry := RepositoryEncoderRegistry{}
			testURL, _ := url.Parse("http://test")
			if err := registry.Register(testURL, new(DefaultRepositoryEncoder)); err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}

			encoder := registry.Lookup(testURL)
			if encoder == nil {
				t.Error("expected encoder to be registered, got nil")
			}
		})
		t.Run("should return the default encoder when one is not registered", func(t *testing.T) {
			t.Parallel()
			testURL, _ := url.Parse("http://test")
			encoder := RepositoryEncoderRegistry{}.Lookup(testURL)
			if encoder == nil {
				t.Error("expected encoder to be registered, got nil")
			}
		})
	})
}
