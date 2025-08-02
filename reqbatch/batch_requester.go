package reqbatch

import (
	"net/http"
	"time"

	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

type BatchRequester interface {
	Do(req []*request.Request) []*response.Response
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
