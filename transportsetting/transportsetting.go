// Package transportsetting provides HTTP transport configuration utilities
// 包 transportsetting 提供 HTTP 传输层配置工具
package transportsetting

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

// TransportSetting manages HTTP transport configuration with thread-safe operations
// TransportSetting 管理 HTTP 传输层配置，提供线程安全操作
type TransportSetting struct {
	transport *http.Transport // HTTP transport instance / HTTP 传输层实例
	mu        sync.Mutex      // Mutex for thread safety / 用于线程安全的互斥锁
}

// SetTLS configures TLS settings with certificate files
// SetTLS 使用证书文件配置 TLS 设置
func (c *TransportSetting) SetTLS(certPath, keyPath, caPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		cert       tls.Certificate
		caCert     []byte
		err        error
		caCertPool *x509.CertPool
	)
	if certPath != "" && keyPath != "" {
		cert, err = tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return fmt.Errorf("failed to load transportsetting cert/key: %w", err)
		}
	}
	if caPath != "" {
		caCert, err = os.ReadFile(caPath)
		if err != nil {
			return fmt.Errorf("failed to read CA cert: %w", err)
		}
		caCertPool = x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return fmt.Errorf("failed to append CA cert to pool")
		}
	}
	tlsConfig := &tls.Config{}
	if cert.Certificate != nil {
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	if caCertPool != nil {
		tlsConfig.RootCAs = caCertPool
	}
	c.transport.TLSClientConfig = tlsConfig
	return nil
}

// SetTransport sets a custom HTTP transport
// SetTransport 设置自定义 HTTP 传输层
func (c *TransportSetting) SetTransport(transport *http.Transport) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport = transport
}

// SetProxy configures proxy settings
// SetProxy 配置代理设置
func (c *TransportSetting) SetProxy(proxies interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if proxies == nil {
		return nil
	}

	switch p := proxies.(type) {
	case string:
		if p == "" {
			c.transport.Proxy = nil
			return nil
		}
		proxyUrl, err := url.Parse(p)
		if err != nil {
			return err
		}
		c.transport.Proxy = http.ProxyURL(proxyUrl)

	case func(r *http.Request) (*url.URL, error):
		c.transport.Proxy = p

	default:
		return fmt.Errorf("invalid proxy type: %T", proxies)
	}

	return nil
}

// SetMaxIdleConns sets the maximum number of idle connections
// SetMaxIdleConns 设置最大空闲连接数
func (c *TransportSetting) SetMaxIdleConns(maxIdleConns int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.MaxIdleConns = maxIdleConns
}

// SetMaxIdleConnsPerHost sets the maximum number of idle connections per host
// SetMaxIdleConnsPerHost 设置每个主机的最大空闲连接数
func (c *TransportSetting) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.MaxIdleConnsPerHost = maxIdleConnsPerHost
}

// SetMaxConnsPerHost sets the maximum number of connections per host
// SetMaxConnsPerHost 设置每个主机的最大连接数
func (c *TransportSetting) SetMaxConnsPerHost(maxConnsPerHost int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.MaxConnsPerHost = maxConnsPerHost
}

// SetIdleConnTimeout sets the idle connection timeout
// SetIdleConnTimeout 设置空闲连接超时时间
func (c *TransportSetting) SetIdleConnTimeout(idleConnTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.IdleConnTimeout = idleConnTimeout
}

// SetTLSHandshakeTimeout sets the TLS handshake timeout
// SetTLSHandshakeTimeout 设置 TLS 握手超时时间
func (c *TransportSetting) SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.TLSHandshakeTimeout = tlsHandshakeTimeout
}

// SetExpectContinueTimeout sets the Expect: 100-continue timeout
// SetExpectContinueTimeout 设置 Expect: 100-continue 超时时间
func (c *TransportSetting) SetExpectContinueTimeout(expectContinueTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.ExpectContinueTimeout = expectContinueTimeout
}

// SetDisableKeepAlives sets whether to disable HTTP Keep-Alive
// SetDisableKeepAlives 设置是否禁用 HTTP Keep-Alive
func (c *TransportSetting) SetDisableKeepAlives(disableKeepAlives bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.DisableKeepAlives = disableKeepAlives
}

// GetTransport returns the configured HTTP transport
// GetTransport 返回配置的 HTTP 传输层
func (c *TransportSetting) GetTransport() *http.Transport {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.transport
}

// NewTransportSetting creates a new transport setting with optional HTTP/2 support
// NewTransportSetting 创建一个新的传输层设置，支持可选的 HTTP/2
func NewTransportSetting(enableHttp2 bool) *TransportSetting {
	result := &TransportSetting{
		transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // Maximum wait time for TCP connection establishment / 建立 TCP 连接的最大等待时间
				KeepAlive: 30 * time.Second, // TCP KeepAlive interval / TCP KeepAlive 的时间间隔
			}).DialContext,
			DisableKeepAlives:     false,            // Whether to disable HTTP Keep-Alive (false means enabled) / 是否禁用 HTTP Keep-Alive（false 表示开启 Keep-Alive）
			MaxIdleConns:          1000,             // Global maximum idle connections / 全局最大空闲连接数
			MaxIdleConnsPerHost:   1000,             // Maximum idle connections per host / 每个主机的最大空闲连接数
			MaxConnsPerHost:       1000,             // Maximum connections per host (including active and idle) / 每个主机的最大连接数（包括正在使用和空闲的）
			IdleConnTimeout:       90 * time.Second, // Idle connections will be closed after 90 seconds / 空闲连接超过 90 秒会被关闭
			TLSHandshakeTimeout:   10 * time.Second, // Maximum TLS handshake wait time / TLS 握手最长等待时间
			ExpectContinueTimeout: 1 * time.Second,  // For HTTP/1.1 Expect: 100-continue mechanism. Client waits for server confirmation before sending body. This field sets wait time, send body directly if no response within 1 second. / 针对 HTTP/1.1 的 Expect: 100-continue 机制。客户端会在发送 body 前等服务器确认。这个字段设定等待时间，1 秒内没响应就直接发 body。
		},
	}

	if enableHttp2 {
		_ = http2.ConfigureTransport(result.transport)
	}
	return result
}
