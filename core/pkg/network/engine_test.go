package network

import (
	"net"
	"testing"
	"time"
)

func TestNetworkEngineScanPortSupportsIPv6(t *testing.T) {
	listener, err := net.Listen("tcp6", "[::1]:0")
	if err != nil {
		t.Skipf("ipv6 loopback unavailable: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	engine, err := NewNetworkEngine(nil)
	if err != nil {
		t.Fatalf("NewNetworkEngine() error = %v", err)
	}

	result := engine.ScanPort("::1", port, 200*time.Millisecond)
	if !result.Open {
		t.Fatalf("expected ipv6 port to be open, got %+v", result)
	}
}
