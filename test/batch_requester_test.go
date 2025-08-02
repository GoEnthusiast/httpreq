package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqbatch"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestBatchGetMethod(t *testing.T) {
	// 创建批量提交请求器
	batchRequester := reqbatch.NewBatchRequester(false)

	// 准备多个请求
	requests := []*request.Request{}
	for i := 0; i < 10; i++ {
		requests = append(requests, &request.Request{
			Method: method.GET,
			URL:    "http://127.0.0.1:9000/testGetNoParams",
		})
	}

	// 执行批量提交请求
	responses := batchRequester.Do(requests)

	// 处理响应
	for _, resp := range responses {
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
