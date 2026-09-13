//go:build android

package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func findEasyTierFFIBesideNativeLibs() string {
	maps, err := os.Open("/proc/self/maps")
	if err != nil {
		return ""
	}
	defer maps.Close()

	scanner := bufio.NewScanner(maps)
	for scanner.Scan() {
		line := scanner.Text()
		path := mapsPath(line)
		if path == "" {
			continue
		}
		base := filepath.Base(path)
		if base != "libcore.so" && base != "libclash.so" {
			continue
		}
		candidate := filepath.Join(filepath.Dir(path), easyTierFFIName)
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() {
			return candidate
		}
	}
	return ""
}
