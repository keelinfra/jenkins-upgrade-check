package uc

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenRetries(t *testing.T) {
	retryWait = 0
	calls, failures := 0, 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls <= failures {
			http.NotFound(w, r)
			return
		}
		io.WriteString(w, "2.580.1")
	}))
	defer srv.Close()
	s := &Source{Base: srv.URL}

	// A mirror that 404s twice and then answers.
	failures = 2
	r, err := s.Open(context.Background(), "stable/latestCore.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(r)
	r.Close()
	if string(b) != "2.580.1" || calls != 3 {
		t.Fatalf("got %q after %d calls", b, calls)
	}

	// One that never answers: give up after three tries.
	calls, failures = 0, 99
	if _, err := s.Open(context.Background(), "stable/latestCore.txt"); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("want a 404, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("want three tries, got %d", calls)
	}
}
