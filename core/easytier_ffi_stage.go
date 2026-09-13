package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/metacubex/mihomo/log"
)

const easyTierFFIName = "libeasytier_ffi.so"

func stageBundledEasyTierFFI(homeDir string) {
	if homeDir == "" {
		return
	}
	src := findBundledEasyTierFFI()
	if src == "" {
		return
	}
	dest := filepath.Join(homeDir, easyTierFFIName)
	if err := copyFileIfChanged(src, dest); err != nil {
		log.Warnln("[EasyTier] stage %s: %v", easyTierFFIName, err)
	}
}

func findBundledEasyTierFFI() string {
	return findEasyTierFFIBesideNativeLibs()
}

func copyFileIfChanged(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	if current, statErr := os.Stat(dest); statErr == nil &&
		current.Size() == info.Size() &&
		current.ModTime().Equal(info.ModTime()) {
		return nil
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), easyTierFFIName+".*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	ok := false
	defer func() {
		if !ok {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := io.Copy(tmp, in); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	if err := os.Chtimes(tmpName, info.ModTime(), info.ModTime()); err != nil {
		return err
	}
	if err := os.Rename(tmpName, dest); err != nil {
		return err
	}
	ok = true
	return nil
}

func mapsPath(line string) string {
	const marker = " /"
	i := strings.Index(line, marker)
	if i < 0 {
		return ""
	}
	return strings.TrimSpace(line[i+1:])
}
