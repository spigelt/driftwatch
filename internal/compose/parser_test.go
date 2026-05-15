package compose_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftwatch/internal/compose"
)

const sampleCompose = `
version: "3.9"
services:
  web:
    image: nginx:1.25
    ports:
      - "80:80"
    environment:
      ENV: production
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: app
      POSTGRES_USER: user
`

func writeTempCompose(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp compose: %v", err)
	}
	return path
}

func TestParseFile_Valid(t *testing.T) {
	path := writeTempCompose(t, sampleCompose)
	cf, err := compose.ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cf.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cf.Services))
	}
	web, ok := cf.Services["web"]
	if !ok {
		t.Fatal("expected service 'web'")
	}
	if web.Image != "nginx:1.25" {
		t.Errorf("expected image nginx:1.25, got %q", web.Image)
	}
}

func TestParseFile_Missing(t *testing.T) {
	_, err := compose.ParseFile("/nonexistent/docker-compose.yml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseFile_NoServices(t *testing.T) {
	path := writeTempCompose(t, "version: \"3\"\n")
	_, err := compose.ParseFile(path)
	if err == nil {
		t.Error("expected error for empty services, got nil")
	}
}
