package anet

import (
	"errors"
	"net"
	"testing"
)

func TestFirstNonEmptyIfacesPrefersPrimary(t *testing.T) {
	primary := []net.Interface{{Name: "wlan0"}}
	fallback := []net.Interface{{Name: "eth0"}}
	got, err := firstNonEmptyIfaces(primary, nil, fallback, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "wlan0" {
		t.Fatalf("got %#v", got)
	}
}

func TestFirstNonEmptyIfacesFallsBackOnError(t *testing.T) {
	fallback := []net.Interface{{Name: "eth0"}}
	got, err := firstNonEmptyIfaces(nil, errors.New("netlink"), fallback, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "eth0" {
		t.Fatalf("got %#v", got)
	}
}

func TestFirstNonEmptyIfacesFallsBackOnEmpty(t *testing.T) {
	fallback := []net.Interface{{Name: "rmnet0"}}
	got, err := firstNonEmptyIfaces(nil, nil, fallback, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "rmnet0" {
		t.Fatalf("got %#v", got)
	}
}

func TestFirstNonEmptyIfacesKeepsPrimaryErrorWhenBothFail(t *testing.T) {
	want := errors.New("netlink")
	_, err := firstNonEmptyIfaces(nil, want, nil, errors.New("stdlib"))
	if err != want {
		t.Fatalf("got %v, want primary error", err)
	}
}

func TestFirstNonEmptyAddrsFallsBackOnEmpty(t *testing.T) {
	fallback := []net.Addr{&net.IPAddr{IP: net.IPv4(192, 168, 1, 2)}}
	got, err := firstNonEmptyAddrs(nil, nil, fallback, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %#v", got)
	}
}

func TestFirstValidIfaceFallsBackOnError(t *testing.T) {
	fallback := &net.Interface{Name: "wlan0"}
	got, err := firstValidIface(nil, errors.New("netlink"), fallback, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Name != "wlan0" {
		t.Fatalf("got %#v", got)
	}
}

func TestChooseIfacesForceAnetFallsBackToStdlib(t *testing.T) {
	std := []net.Interface{{Name: "wlan0"}}
	got, err := chooseIfaces(true, nil, errors.New("netlink"), std, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "wlan0" {
		t.Fatalf("got %#v", got)
	}
}

func TestChooseIfacesForceAnetPrefersNetlink(t *testing.T) {
	netlink := []net.Interface{{Name: "rmnet0"}}
	std := []net.Interface{{Name: "wlan0"}}
	got, err := chooseIfaces(true, netlink, nil, std, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "rmnet0" {
		t.Fatalf("got %#v", got)
	}
}

func TestChooseIfacesForceAnetFallsBackOnEmptyList(t *testing.T) {
	std := []net.Interface{{Name: "wlan0"}}
	got, err := chooseIfaces(true, nil, nil, std, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "wlan0" {
		t.Fatalf("got %#v", got)
	}
}

func TestChooseIfacesWithoutForcePrefersStdlib(t *testing.T) {
	netlink := []net.Interface{{Name: "rmnet0"}}
	std := []net.Interface{{Name: "wlan0"}}
	got, err := chooseIfaces(false, netlink, nil, std, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "wlan0" {
		t.Fatalf("got %#v", got)
	}
}
