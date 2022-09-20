package transport

import (
	"net/url"
	"testing"
)

func TestRepositoryEncoderRegistry(t *testing.T) {
	t.Run("RepositoryEncoderKey", func(t *testing.T) {
		t.Run("NewRepositoryEncoderKey", func(t *testing.T) {
			u, _ := url.Parse("https://api.pro.coinbase.com")
			key := NewRepositoryEncoderKey(u)
			if key.host != "api.pro.coinbase.com" {
				t.Errorf("expected key.host to be %s, got %s", "api.pro.coinbase.com", key.host)
			}
		})

		t.Run("NewDefaultRepositoryEncoderKey", func(t *testing.T) {
			key := NewDefaultRepositoryEncoderKey()
			if key.host != "" {
				t.Errorf("expected key.host to be %s, got %s", "", key.host)
			}
		})
	})

	t.Run("Register", func(t *testing.T) {
		t.Run("should return ErrRepositoryEncoderExists when registering an encoder that already exists",
			func(t *testing.T) {
				registry := RepositoryEncoderRegistry{}
				key := RepositoryEncoderKey{host: "test"}
				registry[key] = new(DefaultRepositoryEncoder)

				u, _ := url.Parse("http://test")
				err := registry.Register(u, new(DefaultRepositoryEncoder))
				if err != ErrRepositoryEncoderExists {
					t.Errorf("expected error %v, got %v", ErrRepositoryEncoderExists, err)
				}
			})

		t.Run("should register an encoder when one does not already exist", func(t *testing.T) {
			registry := RepositoryEncoderRegistry{}
			u, _ := url.Parse("http://test")
			err := registry.Register(u, new(DefaultRepositoryEncoder))
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}

			key := NewRepositoryEncoderKey(u)
			if registry[key] == nil {
				t.Error("expected encoder to be registered, got nil")
			}
		})
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Run("should return the encoder when one is registered", func(t *testing.T) {
			registry := RepositoryEncoderRegistry{}
			u, _ := url.Parse("http://test")
			if err := registry.Register(u, new(DefaultRepositoryEncoder)); err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}

			encoder := registry.Lookup(u)
			if encoder == nil {
				t.Error("expected encoder to be registered, got nil")
			}
		})

		t.Run("should return the default encoder when one is not registered", func(t *testing.T) {
			u, _ := url.Parse("http://test")
			encoder := RepositoryEncoders.Lookup(u)
			if encoder == nil {
				t.Error("expected encoder to be registered, got nil")
			}
		})
	})
}
