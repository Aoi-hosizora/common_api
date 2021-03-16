package service

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib-more/xcharset"
	"golang.org/x/text/encoding"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpService struct{}

func NewHttpService() *HttpService {
	return &HttpService{}
}

var errStatusNotOk = errors.New("service: status code is not 200")

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

func (h *HttpService) HttpGetWithCharset(url string, decode encoding.Encoding, fn func(*http.Request)) ([]byte, *http.Response, error) {
	bs, resp, err := h.HttpGet(url, fn)
	if err != nil {
		return nil, nil, err
	}

	bs, err = xcharset.DecodeBytes(decode, bs)
	if err != nil {
		return nil, nil, err
	}
	return bs, resp, nil
}

func (h *HttpService) HttpPost(url string, body io.Reader, fn func(*http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}
	if fn != nil {
		fn(req)
	}

	return h.DoRequest(req)
}
