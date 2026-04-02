package network

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/lieying/engine/internal/config"
	"golang.org/x/net/proxy"
)

type NetworkManager struct {
	cfg    config.NetworkConfig
	client *http.Client
}

func NewNetworkManager(cfg config.NetworkConfig) (*NetworkManager, error) {
	nm := &NetworkManager{
		cfg: cfg,
	}

	if err := nm.initHTTPClient(); err != nil {
		return nil, err
	}

	return nm, nil
}

func (nm *NetworkManager) initHTTPClient() error {
	transport := &http.Transport{
		DialContext: nm.dialContext,
	}

	switch nm.cfg.Mode {
	case config.NetworkModeOffline:
		transport.DialContext = nm.blockAllDial
	case config.NetworkModeProxy:
		if nm.cfg.Proxy.Enabled {
			proxyURL, err := nm.buildProxyURL()
			if err != nil {
				return err
			}
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	case config.NetworkModeDirect:
	}

	nm.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return nil
}

func (nm *NetworkManager) buildProxyURL() (*url.URL, error) {
	var proxyURLStr string
	switch nm.cfg.Proxy.Type {
	case "socks5":
		if nm.cfg.Proxy.Username != "" {
			proxyURLStr = fmt.Sprintf("socks5://%s:%s@%s:%d",
				nm.cfg.Proxy.Username,
				nm.cfg.Proxy.Password,
				nm.cfg.Proxy.Host,
				nm.cfg.Proxy.Port)
		} else {
			proxyURLStr = fmt.Sprintf("socks5://%s:%d",
				nm.cfg.Proxy.Host,
				nm.cfg.Proxy.Port)
		}
	case "http", "https":
		if nm.cfg.Proxy.Username != "" {
			proxyURLStr = fmt.Sprintf("%s://%s:%s@%s:%d",
				nm.cfg.Proxy.Type,
				nm.cfg.Proxy.Username,
				nm.cfg.Proxy.Password,
				nm.cfg.Proxy.Host,
				nm.cfg.Proxy.Port)
		} else {
			proxyURLStr = fmt.Sprintf("%s://%s:%d",
				nm.cfg.Proxy.Type,
				nm.cfg.Proxy.Host,
				nm.cfg.Proxy.Port)
		}
	default:
		return nil, fmt.Errorf("unsupported proxy type: %s", nm.cfg.Proxy.Type)
	}

	return url.Parse(proxyURLStr)
}

func (nm *NetworkManager) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if !nm.isAllowed(addr) {
		return nil, fmt.Errorf("connection to %s is blocked", addr)
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return dialer.DialContext(ctx, network, addr)
}

func (nm *NetworkManager) blockAllDial(ctx context.Context, network, addr string) (net.Conn, error) {
	return nil, fmt.Errorf("offline mode: all network connections blocked")
}

func (nm *NetworkManager) isAllowed(addr string) bool {
	if len(nm.cfg.BlockedHosts) > 0 {
		for _, blocked := range nm.cfg.BlockedHosts {
			if addr == blocked {
				return false
			}
		}
	}

	if len(nm.cfg.AllowedHosts) > 0 {
		for _, allowed := range nm.cfg.AllowedHosts {
			if addr == allowed {
				return true
			}
		}
		return false
	}

	return true
}

func (nm *NetworkManager) GetClient() *http.Client {
	return nm.client
}

func (nm *NetworkManager) GetMode() config.NetworkMode {
	return nm.cfg.Mode
}

func (nm *NetworkManager) IsOffline() bool {
	return nm.cfg.Mode == config.NetworkModeOffline
}
