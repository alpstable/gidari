package repository

import "testing"

func TestRaw(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		t.Run("should encode the Raw struct into a byte slice", func(t *testing.T) {
			raw := NewRaw("test", []byte("test"))
			_, err := raw.Encode()
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
		})
	})
	t.Run("Decode", func(t *testing.T) {
		t.Run("should decode a byte slice into a Raw struct", func(t *testing.T) {
			raw := NewRaw("test", []byte("test"))

			b, _ := raw.Encode()
			err := raw.Decode(b)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
		})
	})
}
