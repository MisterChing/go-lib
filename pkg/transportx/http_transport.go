package transportx

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

var (
	// DefaultNonKeepAliveTransport 不开启 http-keepalive
	DefaultNonKeepAliveTransport = NewHttpTransport(true, 90*time.Second, 30*time.Second)
	// DefaultFastCloseTransport 开启http-keepalive且空闲tcp最长存活15s
	DefaultFastCloseTransport = NewHttpTransport(false, 15*time.Second, 5*time.Second)
)

func NewHttpTransport(disableKeepAlive bool, idleConnTimeout, probeInterval time.Duration) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	if probeInterval > 0 {
		dialer.KeepAlive = probeInterval
	}
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		DisableKeepAlives:     disableKeepAlive,
		IdleConnTimeout:       90 * time.Second, //在开启http-keepalive下 空闲tcp最长存活时间
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	if idleConnTimeout > 0 {
		tr.IdleConnTimeout = idleConnTimeout
	}
	return tr
}
