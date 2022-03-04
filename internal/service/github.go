package service

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/static"
	"net/http"
	"strings"
)

type GithubService struct {
	httpService *HttpService
}

func NewGithubService() *GithubService {
	return &GithubService{
		httpService: xmodule.MustGetByName(sn.SHttpService).(*HttpService),
	}
}

func (g *GithubService) GetRateLimit(token string) (map[string]interface{}, error) {
	bs, _, err := g.httpService.HttpGet(static.GithubRateLimitApi, func(r *http.Request) {
		r.Header.Add("Authorization", token)
		r.Header.Add("Accept", static.GithubAccept)
	})
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GithubService) RequestApiWithToken(url string, token string) (bs []byte, statusCode int, header http.Header, err error) {
	u := fmt.Sprintf(static.GithubApiPrefix, strings.TrimPrefix(url, "/"))
	bs, resp, err := g.httpService.HttpGetWithoutCheckCode(u, func(r *http.Request) {
		r.Header.Add("Authorization", token)
		r.Header.Add("Accept", static.GithubAccept)
	})
	if err != nil {
		return nil, 0, nil, err
	}
	return bs, resp.StatusCode, resp.Header, nil
}
