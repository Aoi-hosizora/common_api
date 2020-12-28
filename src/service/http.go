package service

import (
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

func (h *HttpService) HttpGet(url string, fn func(*http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	if fn != nil {
		fn(req)
	}

	return h.DoRequest(req)
}
