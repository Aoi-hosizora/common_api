package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/model/obj"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/static"
	"net/http"
	"net/url"
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

func (g *GithubService) GetRepoIssuesByTitle(owner, repo string, page, limit int32, q string, token string) (int32, []*obj.GithubIssueItem, error) {
	qString := fmt.Sprintf("repo:%s/%s is:issue", owner, repo)
	if q != "" {
		qString += fmt.Sprintf(" %s in:title", q)
	}
	qString = url.PathEscape(qString)
	u := fmt.Sprintf(static.GithubIssueSimpleSearchApi, qString, page, limit) // q page per_page
	bs, resp, err := g.httpService.HttpGet(u, func(r *http.Request) {
		r.Header.Set("Authorization", token)
	})
	if err != nil {
		return 0, nil, err
	}
	if resp.StatusCode != 200 {
		return 0, nil, errors.New("response status is not 200 OK")
	}

	r := &obj.GithubIssueSearchResult{}
	err = json.Unmarshal(bs, r)
	if err != nil {
		return 0, nil, err
	}

	total := int32(-1)
	if !r.IncompleteResults {
		total = r.TotalCount
	}
	return total, r.Items, nil
}
