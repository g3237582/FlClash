//go:build amd64

package easytierffi

import (
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

// socketAddrStack packs DataPlaneSocketAddr for SysV MEMORY-class stack passing.
func socketAddrStack(addr SocketAddr) [3]uintptr {
	var words [3]uintptr
	src := unsafe.Slice((*byte)(unsafe.Pointer(&addr)), int(unsafe.Sizeof(addr)))
	dst := unsafe.Slice((*byte)(unsafe.Pointer(&words[0])), int(unsafe.Sizeof(words)))
	copy(dst, src)
	return words
}

func tcpConnectSubmitArgs(session uint64, addr SocketAddr, timeoutMS uint64, out *uint64) []uintptr {
	words := socketAddrStack(addr)
	return []uintptr{
		uintptr(session),
		uintptr(timeoutMS),
		uintptr(unsafe.Pointer(out)),
		0, 0, 0,
		words[0], words[1], words[2],
	}
}

func udpSendSubmitArgs(session, socket uint64, addr SocketAddr, data *byte, length uint32, out *uint64) []uintptr {
	words := socketAddrStack(addr)
	return []uintptr{
		uintptr(session),
		uintptr(socket),
		uintptr(unsafe.Pointer(data)),
		uintptr(length),
		uintptr(unsafe.Pointer(out)),
		0,
		words[0], words[1], words[2],
	}
}

// callTCPConnectSubmit uses SyscallN because the 20-byte addr is passed on the
// stack (SysV MEMORY), not as a register pointer. RegisterLibFunc cannot do that
// on linux. Layout: rdi=session, rsi=timeout, rdx=out; stack=addr.
func (n *Native) callTCPConnectSubmit(session uint64, addr SocketAddr, timeoutMS uint64, out *uint64) int32 {
	args := tcpConnectSubmitArgs(session, addr, timeoutMS, out)
	r1, _, _ := purego.SyscallN(n.tcpConnectSubmitSym, args...)
	runtime.KeepAlive(addr)
	runtime.KeepAlive(out)
	return int32(r1)
}

// callUDPSendSubmit: rdi=session, rsi=socket, rdx=data, rcx=len, r8=out; stack=addr.
func (n *Native) callUDPSendSubmit(session, socket uint64, addr SocketAddr, data *byte, length uint32, out *uint64) int32 {
	args := udpSendSubmitArgs(session, socket, addr, data, length, out)
	r1, _, _ := purego.SyscallN(n.udpSendSubmitSym, args...)
	runtime.KeepAlive(addr)
	runtime.KeepAlive(data)
	runtime.KeepAlive(out)
	return int32(r1)
}
