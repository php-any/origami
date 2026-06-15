package http

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestBufferedWriter_StatusBeforeHeaderRedirect(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	bw.SetStatus(302)
	bw.SetHeader("Location", "/target")
	_, _ = bw.Write(nil)

	if rec.Code != 302 {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/target" {
		t.Fatalf("Location = %q, want /target", got)
	}
}

func TestBufferedWriter_HeaderBeforeStatusRedirect(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	bw.SetHeader("Location", "/target")
	bw.SetStatus(302)
	_, _ = bw.Write([]byte{})

	if rec.Code != 302 {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/target" {
		t.Fatalf("Location = %q, want /target", got)
	}
}

func TestBufferedWriter_Redirect(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	bw.Redirect("/login", 302)

	if rec.Code != 302 {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/login" {
		t.Fatalf("Location = %q, want /login", got)
	}
}

func TestBufferedWriter_NoContent(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	bw.NoContent(204)

	if rec.Code != 204 {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}

func TestBufferedWriter_CommitPending(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	bw.SetStatus(204)
	bw.commitPending()

	if rec.Code != 204 {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}

func TestBufferedWriter_WriteFlushesDefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	_, _ = bw.Write([]byte("ok"))

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body = %q, want ok", rec.Body.String())
	}
}

func TestBufferedWriter_WriteHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)

	if err := bw.WriteHTML([]byte("<p>hi</p>")); err != nil {
		t.Fatal(err)
	}
	if got := rec.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	if rec.Body.String() != "<p>hi</p>" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestBufferedWriter_SendFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(path, []byte("download-me"), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	bw := newBufferedWriter(rec)
	if err := bw.SendFile(path, "greeting.txt"); err != nil {
		t.Fatal(err)
	}

	if got := rec.Header().Get("Content-Disposition"); got != `attachment; filename="greeting.txt"` {
		t.Fatalf("Content-Disposition = %q", got)
	}
	if rec.Body.String() != "download-me" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestNewBufferedWriter_Idempotent(t *testing.T) {
	rec := httptest.NewRecorder()
	bw1 := newBufferedWriter(rec)
	bw2 := newBufferedWriter(bw1)
	if bw1 != bw2 {
		t.Fatal("newBufferedWriter should be idempotent")
	}
}
