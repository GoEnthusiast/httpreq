package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqsingle"
	"github.com/GoEnthusiast/httpreq/types/request"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

// TestSingleGetMethod 简单GET请求
func TestSingleGetMethod(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasParams 带参数的GET请求
func TestSingleGetMethodHasParams(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get?name=GoEnthusiast&age=18",
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasHeaders 带Headers的GET请求
func TestSingleGetMethodHasHeaders(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
		Header: map[string][]string{
			"User-Agent": {"GoEnthusiast"},
		},
	}

	resp := requester.Do(req)
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

// TestSinglePostMethodHasJsonParams 带JSON参数的POST请求
func TestSinglePostMethodHasJsonParams(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.POST,
		URL:    "https://httpbin.org/post",
		Body: map[string]interface{}{
			"name": "GoEnthusiast",
			"age":  18,
		},
		ContentType: method.ContentTypeJSON,
	}

	resp := requester.Do(req)
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

// TestSinglePostMethodHasFormParams 带表单参数的 POST 请求
func TestSinglePostMethodHasFormParams(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.POST,
		URL:    "https://httpbin.org/post",
		Body: map[string]interface{}{
			"name": "GoEnthusiast",
			"age":  18,
		},
		ContentType: method.ContentTypeForm,
	}

	resp := requester.Do(req)
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

// TestSinglePostMethodHasMultipartFormParams 带多部分表单参数的 POST 请求[文件上传]
func TestSinglePostMethodHasMultipartFormParams(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	file, _ := os.Open("test.txt")
	req := &request.Request{
		Method: method.POST,
		URL:    "https://httpbin.org/post",
		Body: map[string]interface{}{
			"name": "GoEnthusiast",
			"age":  18,
			"file": file,
		},
		ContentType: method.ContentTypeMulti,
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasFixedProxyInRequest 在请求体中使用固定代理
func TestSingleGetMethodHasFixedProxyInRequest(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
		Proxy:  "http://221.1.133.210:8080",
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasRandomProxyInRequest 在请求体中使用随机代理
func TestSingleGetMethodHasRandomProxyInRequest(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
		Proxy: func(r *http.Request) (*url.URL, error) {
			// myProxy := getProxy() 自定义实现IP池获取 IP 逻辑
			myProxy := "http://221.1.133.210:8080"
			return url.Parse(myProxy)
		},
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasFixedProxyInRequester 在请求器中设置固定代理
func TestSingleGetMethodHasFixedProxyInRequester(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)
	if err := requester.SetProxy("http://221.1.133.210:8080"); err != nil {
		t.Logf("设置代理错误: %v\n", err)
		return
	}

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasRandomProxyInRequester 在请求器中设置随机代理
func TestSingleGetMethodHasRandomProxyInRequester(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)
	if err := requester.SetProxy(func(r *http.Request) (*url.URL, error) {
		// myProxy := getProxy() 自定义实现IP池获取 IP 逻辑
		myProxy := "http://221.1.133.210:8080"
		return url.Parse(myProxy)
	}); err != nil {
		t.Logf("设置代理错误: %v\n", err)
		return
	}

	req := &request.Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
	}

	resp := requester.Do(req)
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

// TestSingleGetMethodHasTimeout 设置请求超时
func TestSingleGetMethodHasTimeout(t *testing.T) {
	requester := reqsingle.NewSingleRequester(false)

	req := &request.Request{
		Method:  method.GET,
		URL:     "https://httpbin.org/get",
		Timeout: 5 * time.Second,
	}

	resp := requester.Do(req)
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
