package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GoEnthusiast/httpreq/method"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"strings"
)

func BuildRequestBody(contentType method.HTTPContentType, body interface{}) (io.Reader, string, error) {
	if body == nil {
		return nil, string(contentType), nil
	}
	switch contentType {
	case method.ContentTypeJSON:
		// application/json
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, string(contentType), err
		}
		return bytes.NewReader(jsonBytes), string(contentType), nil
	case method.ContentTypeForm:
		// application/x-www-form-urlencoded
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
		// multipart/form-data
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)

		switch v := body.(type) {
		case map[string]interface{}:
			for key, val := range v {
				if file, ok := val.(*os.File); ok {
					fileWriter, err := writer.CreateFormFile(key, file.Name())
					if err != nil {
						return nil, "", err
					}
					_, err = io.Copy(fileWriter, file)
					if err != nil {
						return nil, "", err
					}
				} else {
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
		return buf, writer.FormDataContentType(), nil // 注意 content-type 要由 writer 提供
	case method.ContentTypeText:
		str, ok := body.(string)
		if !ok {
			return nil, string(contentType), fmt.Errorf("body must be string for text/plain")
		}
		return strings.NewReader(str), string(contentType), nil
	default:
		return nil, string(contentType), fmt.Errorf("unsupported content type: %s", contentType)
	}
}
