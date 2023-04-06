package service

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xconstant/headers"
	"io"
	"net/http"
)

type HttpService struct{}

func NewHttpService() *HttpService {
	return &HttpService{}
}

const contentType = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"

func (h *HttpService) HttpGet(url string, fn func(*http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add(headers.ContentType, contentType)
	if fn != nil {
		fn(req)
	}
	return h.doRequest(req, true)
}

func (h *HttpService) HttpGetWithoutCheckCode(url string, fn func(*http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add(headers.ContentType, contentType)
	if fn != nil {
		fn(req)
	}
	return h.doRequest(req, false)
}

func (h *HttpService) HttpGetWithCtx(ctx context.Context, url string, fn func(*http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add(headers.ContentType, contentType)
	if fn != nil {
		fn(req)
	}
	return h.doRequest(req, true)
}

func (h *HttpService) HttpPost(url string, body io.Reader, fn func(*http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add(headers.ContentType, contentType)
	if fn != nil {
		fn(req)
	}
	return h.doRequest(req, true)
}

func (h *HttpService) doRequest(req *http.Request, checkCode bool) ([]byte, *http.Response, error) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if checkCode && resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("service: get non-200 response when requesting `%s`", req.URL.String())
	}

	body := resp.Body
	defer body.Close()
	bs, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}

	return bs, resp, nil
}
