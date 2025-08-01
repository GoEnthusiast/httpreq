package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/GoEnthusiast/httpreq/method"
	"golang.org/x/net/http2"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	transport *http.Transport
	client    *http.Client
}

// SetTLS 设置 TLS
func (c *Client) SetTLS(certPath, keyPath, caPath string) error {
	var (
		cert       tls.Certificate
		caCert     []byte
		err        error
		caCertPool *x509.CertPool
	)
	if certPath != "" && keyPath != "" {
		cert, err = tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return fmt.Errorf("failed to load client cert/key: %w", err)
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
	c.applyTransport()
	return nil
}

// SetTransport 设置 transport
func (c *Client) SetTransport(transport *http.Transport) {
	c.transport = transport
	c.applyTransport()
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// SetProxy 设置代理
func (c *Client) SetProxy(proxies string) error {
	if proxies != "" {
		proxyUrl, err := url.Parse(proxies)
		if err != nil {
			return err
		}
		c.transport.Proxy = http.ProxyURL(proxyUrl)
		c.applyTransport()
	}
	return nil
}

// SetMaxIdleConns 设置最大空闲连接数
func (c *Client) SetMaxIdleConns(maxIdleConns int) {
	c.transport.MaxIdleConns = maxIdleConns
	c.applyTransport()
}

// SetMaxIdleConnsPerHost 设置每个主机（host）允许的最大空闲连接数
func (c *Client) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	c.transport.MaxIdleConnsPerHost = maxIdleConnsPerHost
	c.applyTransport()
}

// SetMaxConnsPerHost 设置每个主机允许的最大连接数
func (c *Client) SetMaxConnsPerHost(maxConnsPerHost int) {
	c.transport.MaxConnsPerHost = maxConnsPerHost
	c.applyTransport()
}

// SetIdleConnTimeout 设置空闲连接超时时间
func (c *Client) SetIdleConnTimeout(idleConnTimeout time.Duration) {
	c.transport.IdleConnTimeout = idleConnTimeout
	c.applyTransport()
}

// SetTLSHandshakeTimeout 设置TLS握手超时时间
func (c *Client) SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration) {
	c.transport.TLSHandshakeTimeout = tlsHandshakeTimeout
	c.applyTransport()
}

// SetExpectContinueTimeout 设置Expect: 100-continue 机制的超时时间
func (c *Client) SetExpectContinueTimeout(expectContinueTimeout time.Duration) {
	c.transport.ExpectContinueTimeout = expectContinueTimeout
	c.applyTransport()
}

// GetClient 获取 http.Client
func (c *Client) GetClient() *http.Client {
	return c.client
}

// GetTransport 获取 http.Transport
func (c *Client) GetTransport() *http.Transport {
	return c.transport
}

// BuildRequestBody 构建请求体
func (c *Client) BuildRequestBody(contentType method.HTTPContentType, body interface{}) (io.Reader, string, error) {
	if body == nil {
		return nil, string(contentType), nil
	}
	switch contentType {
	case method.ContentTypeJSON:
		// application/json
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, string(contentType), err
		}
		return bytes.NewReader(jsonBytes), string(contentType), nil
	case method.ContentTypeForm:
		// application/x-www-form-urlencoded
		values := url.Values{}
		switch v := body.(type) {
		case map[string]string:
			for key, val := range v {
				values.Set(key, val)
			}
		case map[string]interface{}:
			for key, val := range v {
				values.Set(key, fmt.Sprintf("%v", val))
			}
		default:
			return nil, string(contentType), fmt.Errorf("invalid body type for form: %T", body)
		}
		return strings.NewReader(values.Encode()), string(contentType), nil
	case method.ContentTypeMulti:
		// multipart/form-data
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)

		switch v := body.(type) {
		case map[string]interface{}:
			for key, val := range v {
				if file, ok := val.(*os.File); ok {
					fileWriter, err := writer.CreateFormFile(key, file.Name())
					if err != nil {
						return nil, "", err
					}
					_, err = io.Copy(fileWriter, file)
					if err != nil {
						return nil, "", err
					}
				} else {
					_ = writer.WriteField(key, fmt.Sprintf("%v", val))
				}
			}
		default:
			return nil, "", fmt.Errorf("unsupported body type for multipart: %T", body)
		}

		err := writer.Close()
		if err != nil {
			return nil, "", err
		}
		return buf, writer.FormDataContentType(), nil // 注意 content-type 要由 writer 提供
	case method.ContentTypeText:
		str, ok := body.(string)
		if !ok {
			return nil, string(contentType), fmt.Errorf("body must be string for text/plain")
		}
		return strings.NewReader(str), string(contentType), nil
	default:
		return nil, string(contentType), fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func (c *Client) applyTransport() {
	c.client.Transport = c.transport
}

func New(enableHttp2 bool) *Client {
	result := &Client{
		client: &http.Client{},
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
	result.applyTransport()
	return result
}
