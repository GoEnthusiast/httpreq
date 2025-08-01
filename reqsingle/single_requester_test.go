package reqsingle

import (
	"github.com/GoEnthusiast/httpreq/method"
	"os"
	"testing"
)

// TestSingleRequesterDoWithGetNoParams 无参数 GET 请求
func TestSingleRequesterDoWithGetNoParams(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.GET,
		URL:    "http://127.0.0.1:9000/testGetNoParams",
	}
	resp, err := s.Do(req)
	if err != nil {
		t.Error(err)
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
		URL:    "http://127.0.0.1:9000/testGetHasParams",
		Body: map[string]interface{}{
			"name": "GoEnthusiast",
			"age":  18,
		},
		ContentType: method.ContentTypeForm,
	}
	resp, err := s.Do(req)
	if err != nil {
		t.Error(err)
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
	resp, err := s.Do(req)
	if err != nil {
		t.Error(err)
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
	resp, err := s.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(resp.ResponseBody))
}
