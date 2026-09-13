//go:build !no_easytier

package outbound

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/metacubex/mihomo/component/easytierffi"
	C "github.com/metacubex/mihomo/constant"
)

var errEasyTierClosed = errors.New("easytier outbound closed")

// EasyTier is a type:easytier outbound. This FlClash tree uses the FFI data
// plane (ffi-library + instance-name + config/config-file). Embedded
// easytier-go from MetaCubeX #3194 is not compiled in; keep Tailscale as-is.
type EasyTier struct {
	*Base
	option     EasyTierOption
	ctx        context.Context
	cancel     context.CancelFunc
	startOnce  sync.Once
	startErr   error
	ffi        *easytierffi.Native
	ffiLibPath string
	ffiSession *easytierffi.Session
	startedFFI bool
}

type EasyTierOption struct {
	BasicOption
	Name                string   `proxy:"name"`
	NetworkName         string   `proxy:"network-name,omitempty"`
	NetworkSecret       string   `proxy:"network-secret,omitempty"`
	Hostname            string   `proxy:"hostname,omitempty"`
	IPv4                string   `proxy:"ipv4,omitempty"`
	DHCP                bool     `proxy:"dhcp,omitempty"`
	Peers               []string `proxy:"peers,omitempty"`
	Listeners           []string `proxy:"listeners,omitempty"`
	NoListener          *bool    `proxy:"no-listener,omitempty"`
	MappedListeners     []string `proxy:"mapped-listeners,omitempty"`
	ExitNodes           []string `proxy:"exit-nodes,omitempty"`
	ProxyNetworks       []string `proxy:"proxy-networks,omitempty"`
	InstanceName        string   `proxy:"instance-name,omitempty"`
	StateDir            string   `proxy:"state-dir,omitempty"`
	UDP                 bool     `proxy:"udp,omitempty"`
	AcceptDNS           *bool    `proxy:"accept-dns,omitempty"`
	EnableExitNode      *bool    `proxy:"enable-exit-node,omitempty"`
	EnableEncryption    *bool    `proxy:"enable-encryption,omitempty"`
	EncryptionAlgorithm string   `proxy:"encryption-algorithm,omitempty"`
	PrivateMode         *bool    `proxy:"private-mode,omitempty"`
	LatencyFirst        *bool    `proxy:"latency-first,omitempty"`
	DisableP2P          *bool    `proxy:"disable-p2p,omitempty"`
	EnableKCPProxy      *bool    `proxy:"enable-kcp-proxy,omitempty"`
	DisableKCPInput     *bool    `proxy:"disable-kcp-input,omitempty"`
	EnableQUICProxy     *bool    `proxy:"enable-quic-proxy,omitempty"`
	DisableQUICInput    *bool    `proxy:"disable-quic-input,omitempty"`
	MTU                 int      `proxy:"mtu,omitempty"`
	TLDDNSZone          string   `proxy:"tld-dns-zone,omitempty"`
	SecureMode          *bool    `proxy:"secure-mode,omitempty"`
	LocalPrivateKey     string   `proxy:"local-private-key,omitempty"`
	LocalPublicKey      string   `proxy:"local-public-key,omitempty"`
	FFILibrary          string   `proxy:"ffi-library,omitempty"`
	Config              string   `proxy:"config,omitempty"`
	ConfigFile          string   `proxy:"config-file,omitempty"`
}

func NewEasyTier(option EasyTierOption) (*EasyTier, error) {
	if strings.TrimSpace(option.FFILibrary) == "" {
		return nil, errors.New("easytier: ffi-library is required")
	}
	return newEasyTierFFI(option)
}

func (e *EasyTier) DialContext(ctx context.Context, metadata *C.Metadata) (_ C.Conn, err error) {
	return e.dialFFI(ctx, metadata)
}

func (e *EasyTier) ListenPacketContext(ctx context.Context, metadata *C.Metadata) (_ C.PacketConn, err error) {
	return e.listenFFI(ctx, metadata)
}

func (e *EasyTier) ProxyInfo() C.ProxyInfo {
	info := e.Base.ProxyInfo()
	info.DialerProxy = e.option.DialerProxy
	return info
}

func (e *EasyTier) IsL3Protocol(*C.Metadata) bool {
	return true
}

func (e *EasyTier) Close() error {
	return e.closeFFI()
}
