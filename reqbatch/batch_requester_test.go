package reqsingle

import (
	"github.com/GoEnthusiast/httpreq/method"
	"testing"
	"time"
)

// TestBatchRequesterDoWithGetNoParams 无参数 GET 请求
func TestBatchRequesterDoWithGetNoParams(t *testing.T) {
	startTime := time.Now()
	var s BatchRequester
	s = NewBatchRequester(false)
	var reqs []Request
	for i := 0; i < 5; i++ {
		req := &Request{
			Method: method.GET,
			URL:    "http://127.0.0.1:9000/testGetNoParams",
		}
		reqs = append(reqs, *req)
	}
	resps := s.Do(reqs)
	t.Log("resps len:", len(resps))
	t.Log("======= 响应内容 ======")
	for _, resp := range resps {
		if resp.Error != nil {
			t.Error(resp.Error)
		}
		t.Log(string(resp.ResponseBody))
		t.Log("------------------------")
	}
	t.Log("总耗时(秒):", time.Now().Sub(startTime).Seconds())
}
