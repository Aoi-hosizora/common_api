package service

import (
	"context"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpService struct{}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (h *HttpService) DoRequest(req *http.Request) ([]byte, *http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	body := resp.Body
	defer body.Close()

	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}

	return bs, resp, nil
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
	return h.DoRequest(req)
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
	return h.DoRequest(req)
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
	return h.DoRequest(req)
}
