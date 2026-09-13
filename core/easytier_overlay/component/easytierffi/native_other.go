//go:build !amd64

package easytierffi

import (
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

func tcpConnectSubmitArgs(session uint64, addr *SocketAddr, timeoutMS uint64, out *uint64) []uintptr {
	return []uintptr{
		uintptr(session),
		uintptr(unsafe.Pointer(addr)),
		uintptr(timeoutMS),
		uintptr(unsafe.Pointer(out)),
	}
}

func udpSendSubmitArgs(session, socket uint64, addr *SocketAddr, data *byte, length uint32, out *uint64) []uintptr {
	return []uintptr{
		uintptr(session),
		uintptr(socket),
		uintptr(unsafe.Pointer(addr)),
		uintptr(unsafe.Pointer(data)),
		uintptr(length),
		uintptr(unsafe.Pointer(out)),
	}
}

// AAPCS64 (Android/Linux arm64) replaces a MEMORY-class 20-byte struct with a
// pointer to a caller-owned copy. SysV stack-word packing does not apply here.
func (n *Native) callTCPConnectSubmit(session uint64, addr SocketAddr, timeoutMS uint64, out *uint64) int32 {
	args := tcpConnectSubmitArgs(session, &addr, timeoutMS, out)
	r1, _, _ := purego.SyscallN(n.tcpConnectSubmitSym, args...)
	runtime.KeepAlive(addr)
	runtime.KeepAlive(out)
	return int32(r1)
}

func (n *Native) callUDPSendSubmit(session, socket uint64, addr SocketAddr, data *byte, length uint32, out *uint64) int32 {
	args := udpSendSubmitArgs(session, socket, &addr, data, length, out)
	r1, _, _ := purego.SyscallN(n.udpSendSubmitSym, args...)
	runtime.KeepAlive(addr)
	runtime.KeepAlive(data)
	runtime.KeepAlive(out)
	return int32(r1)
}
