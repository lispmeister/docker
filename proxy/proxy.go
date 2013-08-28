package proxy

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	UDPConnTrackTimeout = 90 * time.Second
	UDPBufSize          = 2048
)

// A net.Addr where the IP is split into two fields so you can use it as a key
// in a map:
type connTrackKey struct {
	IPHigh uint64
	IPLow  uint64
	Port   int
}

type connTrackMap map[connTrackKey]*net.UDPConn

func newConnTrackKey(addr *net.UDPAddr) *connTrackKey {
	if len(addr.IP) == net.IPv4len {
		return &connTrackKey{
			IPHigh: 0,
			IPLow:  uint64(binary.BigEndian.Uint32(addr.IP)),
			Port:   addr.Port,
		}
	}
	return &connTrackKey{
		IPHigh: binary.BigEndian.Uint64(addr.IP[:8]),
		IPLow:  binary.BigEndian.Uint64(addr.IP[8:]),
		Port:   addr.Port,
	}
}

type Proxy interface {
	// Start forwarding traffic back and forth the front and back-end
	// addresses.
	Run()
	// Stop forwarding traffic and close both ends of the Proxy.
	Close()
	// Return the address on which the proxy is listening.
	FrontendAddr() net.Addr
	// Return the proxied address.
	BackendAddr() net.Addr
}

func NewProxy(frontendAddr, backendAddr net.Addr) (Proxy, error) {
	switch frontendAddr.(type) {
	case *net.UDPAddr:
		return NewUDPProxy(frontendAddr.(*net.UDPAddr), backendAddr.(*net.UDPAddr))
	case *net.TCPAddr:
		return NewTCPProxy(frontendAddr.(*net.TCPAddr), backendAddr.(*net.TCPAddr))
	default:
		panic(fmt.Errorf("Unsupported protocol"))
	}
}
