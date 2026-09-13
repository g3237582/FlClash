//go:build with_gvisor && !no_tailscale

package outbound

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/metacubex/tailscale/ipn"
)

const (
	tailscaleBackendReadyTimeout = 30 * time.Second
	tailscaleIPWait              = 5 * time.Second
	tailscaleIPPoll              = 100 * time.Millisecond
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

func tailscaleNotifyReady(state *ipn.State, v4, v6 netip.Addr) (ready, running bool) {
	if state != nil && *state == ipn.Running {
		return true, true
	}
	return v4.IsValid() || v6.IsValid(), false
}

func waitBackendReady(ctx context.Context, ready <-chan struct{}, readyErr error, life <-chan struct{}, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = tailscaleBackendReadyTimeout
	}
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	select {
	case <-ready:
		return readyErr
	case <-life:
		return context.Canceled
	case <-waitCtx.Done():
		if err := ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("tailscale: timeout waiting for backend Running or TailscaleIPs")
	}
}

func waitTailscaleBindSrc(ctx context.Context, destIs6 bool, pollIPs func() (netip.Addr, netip.Addr), timeout, interval time.Duration) (netip.Addr, error) {
	if pollIPs == nil {
		return tailscaleBindSrc(netip.Addr{}, netip.Addr{}, destIs6)
	}
	if timeout <= 0 {
		timeout = tailscaleIPWait
	}
	if interval <= 0 {
		interval = tailscaleIPPoll
	}
	v4, v6 := pollIPs()
	if src, err := tailscaleBindSrc(v4, v6, destIs6); err == nil {
		return src, nil
	}
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	timer := time.NewTimer(interval)
	defer timer.Stop()
	for {
		select {
		case <-waitCtx.Done():
			v4, v6 = pollIPs()
			return tailscaleBindSrc(v4, v6, destIs6)
		case <-timer.C:
			v4, v6 = pollIPs()
			if src, err := tailscaleBindSrc(v4, v6, destIs6); err == nil {
				return src, nil
			}
			timer.Reset(interval)
		}
	}
}
