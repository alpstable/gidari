package tools

import (
	"testing"
	"time"
)

func TestLogFormatter(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{}
		if lf.String() != "{}" {
			t.Error("expected empty log formatter to be '{}', got", lf.String())
		}
	})
	t.Run("worker", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{WorkerID: 1}
		if lf.String() != "{w:1}" {
			t.Errorf("expected '{w:1}', got '%s'", lf.String())
		}
	})
	t.Run("worker name", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{WorkerName: "worker"}
		if lf.String() != "{worker:worker}" {
			t.Errorf("expected '{worker:worker}', got '%s'", lf.String())
		}
	})
	t.Run("duration", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{Duration: time.Second}
		if lf.String() != "{d:1s}" {
			t.Errorf("expected '{d:1s}', got '%s'", lf.String())
		}
	})
	t.Run("msg", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{Msg: "hello"}
		if lf.String() != "{m:hello}" {
			t.Errorf("expected '{m:hello }', got '%s'", lf.String())
		}
	})
	t.Run("upserted count", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{UpsertedCount: 1}
		if lf.String() != "{u:1}" {
			t.Errorf("expected '{u:1}', got '%s'", lf.String())
		}
	})
	t.Run("matched count", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{MatchedCount: 1}
		if lf.String() != "{c:1}" {
			t.Errorf("expected '{c:1}', got '%s'", lf.String())
		}
	})
	t.Run("all", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{WorkerID: 1, WorkerName: "worker", Duration: time.Second, Msg: "hello", UpsertedCount: 1}
		if lf.String() != "{w:1, worker:worker, d:1s, u:1, m:hello}" {
			t.Errorf("expected '{w:1, worker:worker, d:1s, u:1, m:hello}', got '%s'", lf.String())
		}
	})
}
