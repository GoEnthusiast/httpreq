package reqsingle

import (
	"github.com/GoEnthusiast/httpreq/method"
	"net/http"
	"net/url"
	"os"
	"testing"
)

// TestSingleRequesterDoWithGetNoParams 无参数 GET 请求
func TestSingleRequesterDoWithGetNoParams(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}

// TestSingleRequesterDoWithGetHasParams 有参数 GET 请求
func TestSingleRequesterDoWithGetHasParams(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.GET,
		URL:    "http://127.0.0.1:9000/testGetHasParams?name=GoEnthusiast&age=18",
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}

// TestSingleRequesterDoWithFixedProxy 使用固定代理
func TestSingleRequesterDoWithFixedProxy(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
		Proxy:  "http://HU27BJ815D8783ID:9F41DA0D9516D5CE@http-dyn.abuyun.com:9020",
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}

// TestSingleRequesterDoWithFixedProxy 使用随机代理
func TestSingleRequesterDoWithRandomProxy(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
		Proxy: func(req *http.Request) (*url.URL, error) {
			// 动态获取代理，这里可以调用自定义方法获取不同的代理 IP
			proxyIP := "HU27BJ815D8783ID:9F41DA0D9516D5CE@http-dyn.abuyun.com:9020"
			return url.Parse("http://" + proxyIP)
		},
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}

// TestSingleRequesterDoWithPostJson POST JSON 请求
func TestSingleRequesterDoWithPostJson(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.POST,
		URL:    "http://127.0.0.1:9000/testPostJson",
		Body: map[string]interface{}{
			"name": "GoEnthusiast",
			"age":  18,
		},
		ContentType: method.ContentTypeJSON,
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}

// TestSingleRequesterDoWithPostForm POST Form 请求
func TestSingleRequesterDoWithPostForm(t *testing.T) {
	file, errFile3 := os.Open("/Users/wangyong/Desktop/template3.xlsx")
	defer file.Close()
	if errFile3 != nil {
		t.Error(errFile3)
		return
	}

	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.POST,
		URL:    "http://127.0.0.1:9000/testPostForm",
		Body: map[string]interface{}{
			"name": "GoEnthusiast",
			"age":  18,
			"file": file,
		},
		ContentType: method.ContentTypeMulti,
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}

// POST application/x-www-form-urlencoded 请求
func TestSingleRequesterDoWithPostFormUrlencoded(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.POST,
		URL:    "http://127.0.0.1:9000/testPostFormUrlEncoded",
		Body: map[string]string{
			"name": "GoEnthusiast",
			"age":  "18",
		},
		ContentType: method.ContentTypeForm,
	}
	resp := s.Do(req)
	if resp.Error != nil {
		t.Error(resp.Error.Error())
		return
	}
	t.Log(string(resp.ResponseBody))
}
