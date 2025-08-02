package response

import (
	"github.com/GoEnthusiast/httpreq/types/request"
	"time"
)

type Response struct {
	Request            *request.Request // 请求体
	ResponseStatusCode int              // 响应状态码
	ResponseBody       []byte           // 响应内容
	Error              error            // 错误
	StartTime          time.Time        // 开始时间
	EndTime            time.Time        // 结束时间
	Duration           float64          // 耗时
}
