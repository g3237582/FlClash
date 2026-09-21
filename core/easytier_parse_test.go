package main

import (
	"testing"

	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/constant"
)

func TestParseEasyTierFFIOutbound(t *testing.T) {
	constant.SetHomeDir(t.TempDir())
	proxy, err := adapter.ParseProxy(map[string]any{
		"name":          "easytier",
		"type":          "easytier",
		"ffi-library":   "./libeasytier_ffi.so",
		"instance-name": "mihomo-easytier",
		"config":        "instance_name = \"mihomo-easytier\"\n",
		"udp":           true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if proxy.Type().String() != "EasyTier" {
		t.Fatalf("type = %s", proxy.Type())
	}
}

func TestParseEasyTierRequiresFFILibrary(t *testing.T) {
	constant.SetHomeDir(t.TempDir())
	_, err := adapter.ParseProxy(map[string]any{
		"name":          "easytier",
		"type":          "easytier",
		"instance-name": "mihomo-easytier",
		"config":        "instance_name = \"mihomo-easytier\"\n",
	})
	if err == nil {
		t.Fatal("expected ffi-library to be required")
	}
}
