// Package reqbatch provides batch HTTP request functionality
// 包 reqbatch 提供批量 HTTP 请求功能
package reqbatch

import (
	"net/http"
	"time"

	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

// BatchRequester defines the interface for batch HTTP request operations
// BatchRequester 定义批量 HTTP 请求操作的接口
type BatchRequester interface {
	// Do executes multiple HTTP requests concurrently and returns all responses
	// Do 并发执行多个 HTTP 请求并返回所有响应
	Do(req []*request.Request) []*response.Response

	// SetTLS configures TLS settings with certificate files
	// SetTLS 使用证书文件配置 TLS 设置
	SetTLS(certPath, keyPath, caPath string) error

	// SetTransport sets a custom HTTP transport
	// SetTransport 设置自定义 HTTP 传输层
	SetTransport(transport *http.Transport)

	// SetProxy configures proxy settings
	// SetProxy 配置代理设置
	SetProxy(proxies interface{}) error

	// SetMaxIdleConns sets the maximum number of idle connections
	// SetMaxIdleConns 设置最大空闲连接数
	SetMaxIdleConns(maxIdleConns int)

	// SetMaxIdleConnsPerHost sets the maximum number of idle connections per host
	// SetMaxIdleConnsPerHost 设置每个主机的最大空闲连接数
	SetMaxIdleConnsPerHost(maxIdleConnsPerHost int)

	// SetMaxConnsPerHost sets the maximum number of connections per host
	// SetMaxConnsPerHost 设置每个主机的最大连接数
	SetMaxConnsPerHost(maxConnsPerHost int)

	// SetIdleConnTimeout sets the idle connection timeout
	// SetIdleConnTimeout 设置空闲连接超时时间
	SetIdleConnTimeout(idleConnTimeout time.Duration)

	// SetTLSHandshakeTimeout sets the TLS handshake timeout
	// SetTLSHandshakeTimeout 设置 TLS 握手超时时间
	SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration)

	// SetExpectContinueTimeout sets the Expect: 100-continue timeout
	// SetExpectContinueTimeout 设置 Expect: 100-continue 超时时间
	SetExpectContinueTimeout(expectContinueTimeout time.Duration)

	// SetDisableKeepAlives sets whether to disable HTTP Keep-Alive
	// SetDisableKeepAlives 设置是否禁用 HTTP Keep-Alive
	SetDisableKeepAlives(disableKeepAlives bool)
}
