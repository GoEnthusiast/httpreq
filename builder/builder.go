// Package builder provides utilities for building HTTP request bodies
// 包 builder 提供构建 HTTP 请求体的工具
package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/GoEnthusiast/httpreq/method"
)

// BuildRequestBody builds an HTTP request body based on content type and data
// BuildRequestBody 根据内容类型和数据构建 HTTP 请求体
// Returns: (io.Reader, string, error) - (request body reader, content type, error)
// 返回: (io.Reader, string, error) - (请求体读取器, 内容类型, 错误)
func BuildRequestBody(contentType method.HTTPContentType, body interface{}) (io.Reader, string, error) {
	if body == nil {
		return nil, string(contentType), nil
	}

	switch contentType {
	case method.ContentTypeJSON:
		// Handle JSON content type
		// 处理 JSON 内容类型
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, string(contentType), err
		}
		return bytes.NewReader(jsonBytes), string(contentType), nil

	case method.ContentTypeForm:
		// Handle form data (application/x-www-form-urlencoded)
		// 处理表单数据 (application/x-www-form-urlencoded)
		values := url.Values{}
		switch v := body.(type) {
		case map[string]string:
			for key, val := range v {
				values.Set(key, val)
			}
		case map[string]interface{}:
			for key, val := range v {
				values.Set(key, fmt.Sprintf("%v", val))
			}
		default:
			return nil, string(contentType), fmt.Errorf("invalid body type for form: %T", body)
		}
		return strings.NewReader(values.Encode()), string(contentType), nil

	case method.ContentTypeMulti:
		// Handle multipart form data (multipart/form-data)
		// 处理多部分表单数据 (multipart/form-data)
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)

		switch v := body.(type) {
		case map[string]interface{}:
			for key, val := range v {
				if file, ok := val.(*os.File); ok {
					// Handle file upload
					// 处理文件上传
					fileWriter, err := writer.CreateFormFile(key, file.Name())
					if err != nil {
						return nil, "", err
					}
					_, err = io.Copy(fileWriter, file)
					if err != nil {
						return nil, "", err
					}
				} else {
					// Handle regular form field
					// 处理普通表单字段
					_ = writer.WriteField(key, fmt.Sprintf("%v", val))
				}
			}
		default:
			return nil, "", fmt.Errorf("unsupported body type for multipart: %T", body)
		}

		err := writer.Close()
		if err != nil {
			return nil, "", err
		}
		// Note: content-type should be provided by the writer
		// 注意: content-type 要由 writer 提供
		return buf, writer.FormDataContentType(), nil

	case method.ContentTypeText:
		// Handle plain text content type
		// 处理纯文本内容类型
		str, ok := body.(string)
		if !ok {
			return nil, string(contentType), fmt.Errorf("body must be string for text/plain")
		}
		return strings.NewReader(str), string(contentType), nil

	default:
		return nil, string(contentType), fmt.Errorf("unsupported content type: %s", contentType)
	}
}
