//go:build !amd64

package easytierffi

import (
	"net"
	"testing"
	"unsafe"
)

func TestTCPConnectSubmitAAPCS64Args(t *testing.T) {
	addr, err := ipv4SocketAddr(net.ParseIP("10.77.0.8"), 443)
	if err != nil {
		t.Fatal(err)
	}
	var out uint64
	args := tcpConnectSubmitArgs(1, &addr, 1000, &out)
	if len(args) != 4 {
		t.Fatalf("arg count = %d, want 4 (x0-x3)", len(args))
	}
	if args[0] != 1 || args[2] != 1000 || args[3] != uintptr(unsafe.Pointer(&out)) {
		t.Fatalf("integer regs = %v", args)
	}
	if args[1] != uintptr(unsafe.Pointer(&addr)) {
		t.Fatalf("addr must be a pointer, got %#x", args[1])
	}
	raw := (*[20]byte)(unsafe.Pointer(args[1]))
	if raw[0] != 4 || raw[1] != 0 {
		t.Fatalf("family bytes = %v", raw[:2])
	}
	if uint16(raw[2])|uint16(raw[3])<<8 != 443 {
		t.Fatalf("port bytes = %v", raw[2:4])
	}
	if raw[4] != 10 || raw[5] != 77 || raw[6] != 0 || raw[7] != 8 {
		t.Fatalf("ipv4 bytes = %v", raw[4:8])
	}
}

func TestUDPSendSubmitAAPCS64Args(t *testing.T) {
	addr, err := ipv4SocketAddr(net.ParseIP("10.77.0.8"), 443)
	if err != nil {
		t.Fatal(err)
	}
	var data byte
	var out uint64
	args := udpSendSubmitArgs(1, 7, &addr, &data, 4, &out)
	if len(args) != 6 {
		t.Fatalf("arg count = %d, want 6 (x0-x5)", len(args))
	}
	if args[2] != uintptr(unsafe.Pointer(&addr)) {
		t.Fatalf("addr must be a pointer, got %#x", args[2])
	}
	if args[0] != 1 || args[1] != 7 || args[3] != uintptr(unsafe.Pointer(&data)) || args[4] != 4 {
		t.Fatalf("integer regs = %v", args)
	}
	if args[5] != uintptr(unsafe.Pointer(&out)) {
		t.Fatalf("out = %#x", args[5])
	}
}
