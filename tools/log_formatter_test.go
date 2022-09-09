package tools

import (
	"testing"
	"time"
)

func TestLogFormatter(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		lf := LogFormatter{}
		if lf.String() != "{}" {
			t.Error("expected empty log formatter to be '{}', got", lf.String())
		}
	})
	t.Run("worker", func(t *testing.T) {
		lf := LogFormatter{WorkerID: 1}
		if lf.String() != "{w:1}" {
			t.Errorf("expected '{w:1}', got '%s'", lf.String())
		}
	})
	t.Run("duration", func(t *testing.T) {
		lf := LogFormatter{Duration: time.Second}
		if lf.String() != "{d:1s}" {
			t.Errorf("expected '{d:1s}', got '%s'", lf.String())
		}
	})
	t.Run("msg", func(t *testing.T) {
		lf := LogFormatter{Msg: "hello"}
		if lf.String() != "{m:hello}" {
			t.Errorf("expected '{m:hello }', got '%s'", lf.String())
		}
	})
	t.Run("all", func(t *testing.T) {
		lf := LogFormatter{WorkerID: 1, Duration: time.Second, Msg: "hello"}
		if lf.String() != "{w:1,d:1s,m:hello}" {
			t.Errorf("expected '{w:1,d:1s,m:hello}', got '%s'", lf.String())
		}
	})
}
