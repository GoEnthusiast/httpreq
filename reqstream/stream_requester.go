package reqstream

import (
	"github.com/GoEnthusiast/httpreq/method"
	"net/http"
	"time"
)

type Request struct {
	Method      method.HTTPMethod      // 请求方法
	URL         string                 // 请求地址
	Header      http.Header            // 请求头
	Body        interface{}            // 请求体
	ContentType method.HTTPContentType // 请求内容类型
	Proxy       interface{}            //  代理
	Timeout     time.Duration          // 请求超时时间
	Meta        map[string]interface{} // 请求元数据
}

type Response struct {
	Request            *Request  // 请求体
	ResponseStatusCode int       // 响应状态码
	ResponseBody       []byte    // 响应内容
	Error              error     // 错误
	StartTime          time.Time // 开始时间
	EndTime            time.Time // 结束时间
	Duration           float64   // 耗时
}

type StreamRequester interface {
	Do(req *Request)
	ResponseCh() <-chan *Response
	SetTLS(certPath, keyPath, caPath string) error
	SetTransport(transport *http.Transport)
	SetProxy(proxies interface{}) error
	SetMaxIdleConns(maxIdleConns int)
	SetMaxIdleConnsPerHost(maxIdleConnsPerHost int)
	SetMaxConnsPerHost(maxConnsPerHost int)
	SetIdleConnTimeout(idleConnTimeout time.Duration)
	SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration)
	SetExpectContinueTimeout(expectContinueTimeout time.Duration)
}
