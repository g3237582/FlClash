//go:build amd64

package easytierffi

import (
	"net"
	"testing"
	"unsafe"
)

func TestSocketAddrStackPacking(t *testing.T) {
	addr, err := ipv4SocketAddr(net.ParseIP("10.77.0.8"), 443)
	if err != nil {
		t.Fatal(err)
	}
	words := socketAddrStack(addr)
	raw := (*[20]byte)(unsafe.Pointer(&words[0]))
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

func TestTCPConnectSubmitSysVArgs(t *testing.T) {
	addr, err := ipv4SocketAddr(net.ParseIP("10.77.0.8"), 443)
	if err != nil {
		t.Fatal(err)
	}
	var out uint64
	args := tcpConnectSubmitArgs(1, addr, 1000, &out)
	if len(args) != 9 {
		t.Fatalf("arg count = %d, want 9 (6 GPRs + 3 stack words)", len(args))
	}
	if args[0] != 1 || args[1] != 1000 || args[2] != uintptr(unsafe.Pointer(&out)) {
		t.Fatalf("integer regs = %v", args[:3])
	}
	if args[3] != 0 || args[4] != 0 || args[5] != 0 {
		t.Fatalf("padded GPRs = %v", args[3:6])
	}
	raw := (*[20]byte)(unsafe.Pointer(&args[6]))
	if raw[0] != 4 || uint16(raw[2])|uint16(raw[3])<<8 != 443 {
		t.Fatalf("stack addr = %v", raw[:8])
	}
}

func TestUDPSendSubmitSysVArgs(t *testing.T) {
	addr, err := ipv4SocketAddr(net.ParseIP("10.77.0.8"), 443)
	if err != nil {
		t.Fatal(err)
	}
	var data byte
	var out uint64
	args := udpSendSubmitArgs(1, 7, addr, &data, 4, &out)
	if len(args) != 9 {
		t.Fatalf("arg count = %d, want 9", len(args))
	}
	if args[0] != 1 || args[1] != 7 || args[2] != uintptr(unsafe.Pointer(&data)) || args[3] != 4 {
		t.Fatalf("integer regs = %v", args[:5])
	}
	if args[4] != uintptr(unsafe.Pointer(&out)) || args[5] != 0 {
		t.Fatalf("out/pad = %v", args[4:6])
	}
}
