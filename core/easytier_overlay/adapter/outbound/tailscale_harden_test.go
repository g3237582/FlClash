//go:build with_gvisor && !no_tailscale

package outbound

import (
	"context"
	"errors"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/metacubex/tailscale/ipn"
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

func TestTailscaleNotifyReadyWaitsPastNeedsLogin(t *testing.T) {
	v4 := netip.MustParseAddr("100.64.0.1")
	v6 := netip.MustParseAddr("fd7a:115c:a1e0::1")
	needsLogin := ipn.NeedsLogin
	starting := ipn.Starting
	running := ipn.Running
	noState := ipn.NoState

	cases := []struct {
		name    string
		state   *ipn.State
		v4, v6  netip.Addr
		ready   bool
		running bool
	}{
		{name: "nil-state-no-ips", ready: false, running: false},
		{name: "nostate-no-ips", state: &noState, ready: false, running: false},
		{name: "needslogin-no-ips", state: &needsLogin, ready: false, running: false},
		{name: "starting-no-ips", state: &starting, ready: false, running: false},
		{name: "needslogin-with-v4", state: &needsLogin, v4: v4, ready: true, running: false},
		{name: "starting-with-v6", state: &starting, v6: v6, ready: true, running: false},
		{name: "nil-state-with-v4", v4: v4, ready: true, running: false},
		{name: "running-no-ips", state: &running, ready: true, running: true},
		{name: "running-with-ips", state: &running, v4: v4, v6: v6, ready: true, running: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ready, running := tailscaleNotifyReady(tc.state, tc.v4, tc.v6)
			if ready != tc.ready || running != tc.running {
				t.Fatalf("ready=%v running=%v, want ready=%v running=%v", ready, running, tc.ready, tc.running)
			}
		})
	}
}

func TestWaitBackendReadyUnblocksOnRunning(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	if err := waitBackendReady(context.Background(), ch, nil, nil, time.Second); err != nil {
		t.Fatal(err)
	}
}

func TestWaitBackendReadyKeepsWaitingWithoutIPs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := waitBackendReady(ctx, make(chan struct{}), nil, nil, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v, want context deadline", err)
	}
}

func TestWaitBackendReadyReportsWatchError(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	want := errors.New("watch")
	if err := waitBackendReady(context.Background(), ch, want, nil, time.Second); err != want {
		t.Fatalf("got %v, want watch error", err)
	}
}

func TestWaitBackendReadyDefaultTimeoutMessage(t *testing.T) {
	err := waitBackendReady(context.Background(), make(chan struct{}), nil, nil, 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "Running or TailscaleIPs") {
		t.Fatalf("got %v", err)
	}
}

func TestWaitTailscaleBindSrcImmediate(t *testing.T) {
	v4 := netip.MustParseAddr("100.64.0.1")
	src, err := waitTailscaleBindSrc(context.Background(), false, func() (netip.Addr, netip.Addr) {
		return v4, netip.Addr{}
	}, time.Second, 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if src != v4 {
		t.Fatalf("src = %v, want %v", src, v4)
	}
}

func TestWaitTailscaleBindSrcPollsUntilAssigned(t *testing.T) {
	v4 := netip.MustParseAddr("100.64.0.2")
	var n atomic.Int32
	src, err := waitTailscaleBindSrc(context.Background(), false, func() (netip.Addr, netip.Addr) {
		if n.Add(1) < 3 {
			return netip.Addr{}, netip.Addr{}
		}
		return v4, netip.Addr{}
	}, time.Second, 5*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if src != v4 {
		t.Fatalf("src = %v, want %v", src, v4)
	}
	if n.Load() < 3 {
		t.Fatalf("polls = %d, want at least 3", n.Load())
	}
}

func TestWaitTailscaleBindSrcTimeoutStillEmpty(t *testing.T) {
	_, err := waitTailscaleBindSrc(context.Background(), false, func() (netip.Addr, netip.Addr) {
		return netip.Addr{}, netip.Addr{}
	}, 20*time.Millisecond, 5*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "IPv4") {
		t.Fatalf("got %v", err)
	}
}

func TestWaitTailscaleBindSrcWaitsForDestFamily(t *testing.T) {
	v4 := netip.MustParseAddr("100.64.0.3")
	_, err := waitTailscaleBindSrc(context.Background(), true, func() (netip.Addr, netip.Addr) {
		return v4, netip.Addr{}
	}, 20*time.Millisecond, 5*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "IPv6") {
		t.Fatalf("got %v, want IPv6 miss despite valid v4", err)
	}
}
