package reqstream

import (
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
	"net/http"
	"time"
)

type StreamRequester interface {
	Do(req *request.Request)
	ResponseCh() <-chan *response.Response
	SetTLS(certPath, keyPath, caPath string) error
	SetTransport(transport *http.Transport)
	SetProxy(proxies interface{}) error
	SetMaxIdleConns(maxIdleConns int)
	SetMaxIdleConnsPerHost(maxIdleConnsPerHost int)
	SetMaxConnsPerHost(maxConnsPerHost int)
	SetIdleConnTimeout(idleConnTimeout time.Duration)
	SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration)
	SetExpectContinueTimeout(expectContinueTimeout time.Duration)
	SetDisableKeepAlives(disableKeepAlives bool)
}
