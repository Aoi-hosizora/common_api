package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string, header *http.Header, cb func(string)) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = *header

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	cb(fmt.Sprintf("%d: %s", resp.StatusCode, url))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}

	body := resp.Body
	defer body.Close()

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
