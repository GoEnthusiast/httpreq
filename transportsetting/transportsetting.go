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

type TransportSetting struct {
	transport *http.Transport
	mu        sync.Mutex
}

// SetTLS 设置 TLS
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

// SetTransport 设置自定义 transport
func (c *TransportSetting) SetTransport(transport *http.Transport) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport = transport
}

// SetProxy 设置代理
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

// SetMaxIdleConns 设置最大空闲连接数
func (c *TransportSetting) SetMaxIdleConns(maxIdleConns int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.MaxIdleConns = maxIdleConns
}

// SetMaxIdleConnsPerHost 设置每个主机（host）允许的最大空闲连接数
func (c *TransportSetting) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.MaxIdleConnsPerHost = maxIdleConnsPerHost
}

// SetMaxConnsPerHost 设置每个主机允许的最大连接数
func (c *TransportSetting) SetMaxConnsPerHost(maxConnsPerHost int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.MaxConnsPerHost = maxConnsPerHost
}

// SetIdleConnTimeout 设置空闲连接超时时间
func (c *TransportSetting) SetIdleConnTimeout(idleConnTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.IdleConnTimeout = idleConnTimeout
}

// SetTLSHandshakeTimeout 设置TLS握手超时时间
func (c *TransportSetting) SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.TLSHandshakeTimeout = tlsHandshakeTimeout
}

// SetExpectContinueTimeout 设置Expect: 100-continue 机制的超时时间
func (c *TransportSetting) SetExpectContinueTimeout(expectContinueTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.ExpectContinueTimeout = expectContinueTimeout
}

// SetDisableKeepAlives 设置是否禁用 HTTP Keep-Alive
func (c *TransportSetting) SetDisableKeepAlives(disableKeepAlives bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.transport.DisableKeepAlives = disableKeepAlives
}

// GetTransport 获取 http.Transport
func (c *TransportSetting) GetTransport() *http.Transport {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.transport
}

func NewTransportSetting(enableHttp2 bool) *TransportSetting {
	result := &TransportSetting{
		transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // 建立 TCP 连接的最大等待时间
				KeepAlive: 30 * time.Second, // TCP KeepAlive 的时间间隔
			}).DialContext,
			DisableKeepAlives:     false,            // 是否禁用 HTTP Keep-Alive（false 表示开启 Keep-Alive）
			MaxIdleConns:          1000,             // 全局最大空闲连接数
			MaxIdleConnsPerHost:   1000,             // 每个主机（host）允许的最大空闲连接数
			MaxConnsPerHost:       1000,             // 每个主机允许的最大连接数（包括正在使用和空闲的）
			IdleConnTimeout:       90 * time.Second, // 空闲连接超过 90 秒会被关闭
			TLSHandshakeTimeout:   10 * time.Second, // TLS 握手最长等待时间
			ExpectContinueTimeout: 1 * time.Second,  // 针对 HTTP/1.1 的 Expect: 100-continue 机制。客户端会在发送 body 前等服务器确认。这个字段设定等待时间，1 秒内没响应就直接发 body。
		},
	}

	if enableHttp2 {
		_ = http2.ConfigureTransport(result.transport)
	}
	return result
}
