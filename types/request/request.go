package request

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
