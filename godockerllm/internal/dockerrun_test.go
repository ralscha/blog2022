package internal

import (
	"archive/tar"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSourceArchivePreservesShellCharacters(t *testing.T) {
	code := "package main\n// don't expand $(touch /tmp/unwanted) or `commands`\nfunc main() {}\n"
	archive, err := sourceArchive(code)
	if err != nil {
		t.Fatal(err)
	}
	r := tar.NewReader(archive)
	if _, err := r.Next(); err != nil {
		t.Fatal(err)
	}
	header, err := r.Next()
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != "code/app.go" || string(data) != code {
		t.Fatalf("source was altered: %s %q", header.Name, data)
	}
}

func TestDeadlineStillForcesContainerCleanup(t *testing.T) {
	removed := make(chan bool, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/_ping"):
			w.Header().Set("API-Version", "1.51")
			w.WriteHeader(http.StatusOK)
		case strings.HasSuffix(r.URL.Path, "/images/create"):
			io.WriteString(w, "{}\n")
		case strings.HasSuffix(r.URL.Path, "/containers/create"):
			io.WriteString(w, `{"Id":"test-container","Warnings":[]}`)
		case strings.HasSuffix(r.URL.Path, "/archive"):
			io.Copy(io.Discard, r.Body)
			w.WriteHeader(http.StatusOK)
		case strings.HasSuffix(r.URL.Path, "/start"):
			w.WriteHeader(http.StatusNoContent)
		case strings.HasSuffix(r.URL.Path, "/wait"):
			<-r.Context().Done()
		case r.Method == http.MethodDelete:
			removed <- r.URL.Query().Get("force") == "1"
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected Docker request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	t.Setenv("DOCKER_HOST", "tcp://"+strings.TrimPrefix(server.URL, "http://"))
	t.Setenv("DOCKER_API_VERSION", "1.51")
	t.Setenv("DOCKER_TLS_VERIFY", "")
	_, _, err := RunCodeInDocker("package main\nfunc main() {}", 300*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline, got %v", err)
	}
	select {
	case force := <-removed:
		if !force {
			t.Fatal("cleanup did not force removal")
		}
	case <-time.After(time.Second):
		t.Fatal("container was not removed after deadline")
	}
}
