// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	t.Run("host", func(t *testing.T) {
		t.Parallel()
		lf := LogFormatter{Host: "localhost:8080"}
		if lf.String() != "{host:localhost:8080}" {
			t.Errorf("expected '{host:localhost:8080}', got '%s'", lf.String())
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
		lf := LogFormatter{
			WorkerID: 1, WorkerName: "worker", Duration: time.Second, Host: "localhost",
			Msg: "hello", UpsertedCount: 1,
		}
		if lf.String() != "{w:1, worker:worker, d:1s, host:localhost, u:1, m:hello}" {
			t.Errorf("expected '{w:1, worker:worker, d:1s, host:localhost, u:1, m:hello}', got '%s'", lf.String())
		}
	})
}
