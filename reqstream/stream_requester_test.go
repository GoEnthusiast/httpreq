package reqstream

import (
	"fmt"
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
	"time"
)

// TestStreamRequesterDoWithGetNoParams 流式提交请求
func TestStreamRequesterDoWithGetNoParams(t *testing.T) {
	var s StreamRequester
	s = NewStreamRequester(false, 5)

	go func() {
		for {
			req := &request.Request{
				Method: method.GET,
				URL:    "http://127.0.0.1:9000/testGetNoParams",
			}
			s.Do(req)
			fmt.Println("req:", time.Now().UnixNano())
		}
	}()

	for {
		resp := <-s.ResponseCh()
		if resp.Error != nil {
			t.Error(resp.Error)
		}
		t.Log(string(resp.ResponseBody))
		t.Log("------------------------")
	}
}
