package reqsingle

import (
	"encoding/json"
	"github.com/GoEnthusiast/httpreq/method"
	"testing"
)

func TestSingleRequesterDo(t *testing.T) {
	var s SingleRequester
	s = NewSingleRequester(false)
	req := &Request{
		Method: method.GET,
		URL:    "https://httpbin.org/get",
	}
	resp, err := s.Do(req)
	if err != nil {
		t.Error(err)
	}

	respBytes, _ := json.Marshal(resp)
	t.Log(string(respBytes))
	t.Log(string(resp.ResponseBody))
}
