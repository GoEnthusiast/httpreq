package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqstream"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestStreamGetMethod(t *testing.T) {
	// 创建流式提交请求器，设置并发数为 5
	streamRequester := reqstream.NewStreamRequester(false, 5)

	// 启动请求发送协程
	go func() {
		for i := 0; i < 20; i++ {
			req := &request.Request{
				Method: method.GET,
				URL:    "http://127.0.0.1:9000/testGetNoParams",
			}
			streamRequester.Do(req)
			t.Log("发送请求")
		}
	}()

	// 监听响应
	for i := 0; i < 20; i++ {
		resp := <-streamRequester.ResponseCh()
		if resp.Error != nil {
			t.Logf("请求错误: %v\n", resp.Error)
			return
		}

		t.Logf("状态码: %d\n", resp.ResponseStatusCode)
		t.Logf("响应内容: %s\n", string(resp.ResponseBody))
		t.Logf("请求开始时间: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
		t.Logf("请求结束时间: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
		t.Logf("请求耗时: %.2fs\n", resp.Duration)
	}
}
