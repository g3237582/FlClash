//go:build with_gvisor && !no_tailscale

package outbound

import (
	"net/netip"
	"strings"
	"testing"
)

func TestTailscaleBindSrcRejectsInvalidIPv4(t *testing.T) {
	_, err := tailscaleBindSrc(netip.Addr{}, netip.MustParseAddr("fd7a::1"), false)
	if err == nil {
		t.Fatal("expected error for invalid IPv4 bind src")
	}
	if !strings.Contains(err.Error(), "IPv4") {
		t.Fatalf("error should name IPv4: %v", err)
	}
}

func TestTailscaleBindSrcRejectsInvalidIPv6(t *testing.T) {
	_, err := tailscaleBindSrc(netip.MustParseAddr("100.64.0.1"), netip.Addr{}, true)
	if err == nil {
		t.Fatal("expected error for invalid IPv6 bind src")
	}
	if !strings.Contains(err.Error(), "IPv6") {
		t.Fatalf("error should name IPv6: %v", err)
	}
}

func TestTailscaleBindSrcUsesMatchingFamily(t *testing.T) {
	v4 := netip.MustParseAddr("100.64.0.1")
	v6 := netip.MustParseAddr("fd7a:115c:a1e0::1")

	src, err := tailscaleBindSrc(v4, v6, false)
	if err != nil {
		t.Fatal(err)
	}
	if src != v4 {
		t.Fatalf("src = %v, want %v", src, v4)
	}

	src, err = tailscaleBindSrc(v4, v6, true)
	if err != nil {
		t.Fatal(err)
	}
	if src != v6 {
		t.Fatalf("src = %v, want %v", src, v6)
	}
}

func TestRecoverIntoConvertsPanic(t *testing.T) {
	err := func() (err error) {
		defer recoverInto(&err)
		panic("dial")
	}()
	if err == nil {
		t.Fatal("expected recovered error")
	}
	if !strings.Contains(err.Error(), "tailscale panic: dial") {
		t.Fatalf("got %v", err)
	}
}

func TestRecoverIntoKeepsExistingError(t *testing.T) {
	err := func() (err error) {
		defer recoverInto(&err)
		return errTailscaleNoAddr("IPv4")
	}()
	if err == nil || !strings.Contains(err.Error(), "IPv4") {
		t.Fatalf("got %v", err)
	}
}
