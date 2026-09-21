package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCopyFileIfChanged(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.so")
	dest := filepath.Join(dir, easyTierFFIName)
	payload := []byte("easytier-ffi")
	if err := os.WriteFile(src, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	mtime := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	if err := os.Chtimes(src, mtime, mtime); err != nil {
		t.Fatal(err)
	}
	if err := copyFileIfChanged(src, dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Fatalf("dest = %q", got)
	}
	if err := copyFileIfChanged(src, dest); err != nil {
		t.Fatal(err)
	}
}

func TestMapsPath(t *testing.T) {
	got := mapsPath("7a123000-7a124000 r-xp 00000000 fd:00 1 /data/app/lib/arm64/libcore.so")
	if got != "/data/app/lib/arm64/libcore.so" {
		t.Fatalf("mapsPath = %q", got)
	}
	if mapsPath("anon") != "" {
		t.Fatal("expected empty path")
	}
}
