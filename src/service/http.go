package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpService struct{}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (h *HttpService) HttpRequest(method string, url string, body io.Reader, header *http.Header, cb func(string)) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = *header
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if cb != nil {
		cb(fmt.Sprintf("%d: %s", resp.StatusCode, url))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(resp.Status)
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (h *HttpService) HttpGet(url string, header *http.Header, cb func(string)) ([]byte, error) {
	return h.HttpRequest("GET", url, nil, header, cb)
}
