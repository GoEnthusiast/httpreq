package reqsingle

import (
	"fmt"
	"github.com/GoEnthusiast/httpreq/client"
	"io"
	"net/http"
	"sync"
	"time"
)

type BatchRequesterImpl struct {
	Client *client.Client
}

func (s *BatchRequesterImpl) GetClient() *client.Client {
	return s.Client
}

func (s *BatchRequesterImpl) Do(reqs []Request) []Response {
	var (
		wg        sync.WaitGroup
		mutex     sync.Mutex
		responses []Response
	)

	for i := range reqs {
		req := reqs[i] // 避免闭包引用同一个变量
		wg.Add(1)

		go func() {
			defer wg.Done()

			startTime := time.Now()
			// 设置请求超时
			s.Client.SetTimeout(req.Timeout)

			// 设置代理 IP
			if err := s.Client.SetProxy(req.Proxy); err != nil {
				mutex.Lock()
				responses = append(responses, Response{
					Request:   &req,
					Error:     fmt.Errorf("SetProxy failed for %s: %w", req.URL, err),
					StartTime: startTime,
					EndTime:   time.Now(),
					Duration:  time.Now().Sub(startTime).Seconds(),
				})
				mutex.Unlock()
				return
			}

			// 构造请求体
			body, contentType, err := s.Client.BuildRequestBody(req.ContentType, req.Body)
			if err != nil {
				mutex.Lock()
				responses = append(responses, Response{
					Request:   &req,
					Error:     fmt.Errorf("BuildRequestBody failed for %s: %w", req.URL, err),
					StartTime: startTime,
					EndTime:   time.Now(),
					Duration:  time.Now().Sub(startTime).Seconds(),
				})
				mutex.Unlock()
				return
			}

			httpReq, err := http.NewRequest(string(req.Method), req.URL, body)
			if err != nil {
				mutex.Lock()
				responses = append(responses, Response{
					Request:   &req,
					Error:     fmt.Errorf("NewRequest failed for %s: %w", req.URL, err),
					StartTime: startTime,
					EndTime:   time.Now(),
					Duration:  time.Now().Sub(startTime).Seconds(),
				})
				mutex.Unlock()
				return
			}

			if req.Header != nil {
				httpReq.Header = req.Header
			} else {
				httpReq.Header = make(http.Header)
			}
			if contentType != "" {
				httpReq.Header.Set("Content-Type", contentType)
			}

			// 执行请求
			httpResp, err := s.Client.GetClient().Do(httpReq)
			if err != nil {
				mutex.Lock()
				responses = append(responses, Response{
					Request:   &req,
					Error:     fmt.Errorf("Do failed for %s: %w", req.URL, err),
					StartTime: startTime,
					EndTime:   time.Now(),
					Duration:  time.Now().Sub(startTime).Seconds(),
				})
				mutex.Unlock()
				return
			}
			defer httpResp.Body.Close()

			respBody, err := io.ReadAll(httpResp.Body)
			if err != nil {
				mutex.Lock()
				responses = append(responses, Response{
					Request:   &req,
					Error:     fmt.Errorf("ReadAll failed for %s: %w", req.URL, err),
					StartTime: startTime,
					EndTime:   time.Now(),
					Duration:  time.Now().Sub(startTime).Seconds(),
				})
				mutex.Unlock()
				return
			}

			endTime := time.Now()

			resp := Response{
				Request:            &req,
				ResponseStatusCode: httpResp.StatusCode,
				ResponseBody:       respBody,
				StartTime:          startTime,
				EndTime:            endTime,
				Duration:           endTime.Sub(startTime).Seconds(),
			}

			mutex.Lock()
			responses = append(responses, resp)
			mutex.Unlock()
		}()
	}

	wg.Wait()

	return responses
}

func NewBatchRequester(enableHttp2 bool) BatchRequester {
	c := client.New(enableHttp2)
	return &BatchRequesterImpl{
		Client: c,
	}
}
