package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xerror"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/common_api/internal/model/object"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/static"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type GithubService struct {
	httpService *HttpService
}

func NewGithubService() *GithubService {
	return &GithubService{
		httpService: xmodule.MustGetByName(sn.SHttpService).(*HttpService),
	}
}

func (g *GithubService) GetRateLimit(token string) (map[string]any, error) {
	bs, _, err := g.httpService.HttpGet(static.GithubRateLimitApi, func(r *http.Request) {
		r.Header.Add("Authorization", token)
		r.Header.Add("Accept", static.GithubAccept)
	})
	if err != nil {
		return nil, err
	}

	data := make(map[string]any)
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

func (g *GithubService) GetRepoIssuesByTitle(owner, repo string, page, limit uint32, q string, token string) (uint32, []*object.GithubIssueItem, error) {
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

	r := &object.GithubIssueSearchResult{}
	err = json.Unmarshal(bs, r)
	if err != nil {
		return 0, nil, err
	}

	total := uint32(len(r.Items))
	if !r.IncompleteResults {
		total = uint32(r.TotalCount)
	}
	return total, r.Items, nil
}

func (g *GithubService) GetAoiHosizoraUserProfile(token string) (map[string]any, error) {
	const userApiUrl = "users/Aoi-hosizora"
	var bs1, bs2 []byte
	var code1, code2 int
	var err1, err2 error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		bs1, code1, _, err1 = g.RequestApiWithToken(userApiUrl, "")
	}()
	go func() {
		defer wg.Done()
		bs2, code2, _, err2 = g.RequestApiWithToken(userApiUrl, token)
	}()
	wg.Wait()
	if err1 != nil || err2 != nil {
		return nil, xerror.Combine(err1, err2)
	}
	if code1 != 200 || code2 != 200 {
		return nil, errors.New("failed to make request to users api")
	}

	obj1 := make(map[string]any) // without token
	obj2 := make(map[string]any) // with token
	if json.Unmarshal(bs1, &obj1) != nil || json.Unmarshal(bs2, &obj2) != nil {
		return nil, errors.New("failed to unmarshal response to json")
	}

	newFieldKeys := []string{"private_gists", "owned_private_repos", "total_private_repos"}
	for _, key := range newFieldKeys {
		if field, ok := obj2[key]; ok {
			obj1[key] = field
		}
	}
	needToCombineKeys := [][3]string{{"public_repos", "owned_private_repos", "total_repos"}, {"public_gists", "private_gists", "total_gists"}}
	for _, keys := range needToCombineKeys {
		key1, key2, newKey := keys[0], keys[1], keys[2]
		if privateVal, ok := obj2[key2]; ok {
			if publicVal, ok := obj1[key1]; ok {
				if private, ok := privateVal.(float64); ok {
					if public, ok := publicVal.(float64); ok {
						obj1[newKey] = int32(private + public)
					}
				}
			}
		}
	}
	return obj1, nil
}

func (g *GithubService) GetAoiHosizoraContribution(token string) (map[string]int32, error) {
	username := "Aoi-hosizora"
	from := 2019
	to := time.Now().Year()
	token = "Bearer " + strings.TrimPrefix(token, "Token ")
	template := `{"query":"query {\n  user(login: \"%s\") {\n    name\n    contributionsCollection(from: \"%d-01-01T00:00:00Z\", to: \"%d-12-31T23:59:59Z\") {\n      contributionCalendar {\n        totalContributions\n      }\n    }\n  }\n}"}`

	/*
		curl
		-H "Authorization: bearer xxx"
		-X POST
		-d '{"query":"query {\n  user(login: \"Aoi-hosizora\") {\n    name\n    contributionsCollection(from: \"2020-01-01T00:00:00Z\", to: \"2020-12-31T23:59:59Z\") {\n      contributionCalendar {\n        colors\n        totalContributions\n      }\n    }\n  }\n}"}'
		https://api.github.com/graphql
	*/

	mutex := sync.Mutex{}
	out := make(map[string]int32)

	eg := xerror.NewErrorGroup(context.Background())
	for year := from; year <= to; year++ {
		year := year
		eg.Go(func(c context.Context) error {
			body := fmt.Sprintf(template, username, year, year)
			bs, _, err := g.httpService.HttpPost(static.GithubGraphqlApi, strings.NewReader(body), func(r *http.Request) {
				r.Header.Add("Authorization", token)
				r.Header.Add("Accept", static.GithubAccept)
			})
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`"totalContributions":(\d*)`)
			matches := re.FindStringSubmatch(string(bs))
			if len(matches) < 2 {
				return errors.New("failed to query totalContributions")
			}
			contribution, err := xnumber.Atoi32(matches[1])
			if err != nil {
				return err
			}

			mutex.Lock()
			out[xnumber.Itoa(year)] = contribution
			mutex.Unlock()
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	return out, nil
}
