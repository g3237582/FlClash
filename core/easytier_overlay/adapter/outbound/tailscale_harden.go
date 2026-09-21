//go:build with_gvisor && !no_tailscale

package outbound

import (
	"fmt"
	"net/netip"
)

func errTailscaleNoAddr(family string) error {
	return fmt.Errorf("tailscale: no valid %s address (backend not ready)", family)
}

func tailscaleBindSrc(v4, v6 netip.Addr, destIs6 bool) (netip.Addr, error) {
	src := v4
	family := "IPv4"
	if destIs6 {
		src = v6
		family = "IPv6"
	}
	if !src.IsValid() {
		return netip.Addr{}, errTailscaleNoAddr(family)
	}
	return src, nil
}

func recoverInto(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("tailscale panic: %v", r)
	}
}
