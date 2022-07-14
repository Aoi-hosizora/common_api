package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/model/biz"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/static"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

func (g *GithubService) GetRateLimit(auth string) (map[string]interface{}, error) {
	bs, _, err := g.httpService.HttpGet(static.GithubRateLimitApi, func(r *http.Request) {
		r.Header.Add("Authorization", auth)
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

func (g *GithubService) GetRepoIssuesByTitle(owner, repo string, page, limit int32, q string, auth string) (int32, []*biz.GithubIssueItem, error) {
	qString := fmt.Sprintf("repo:%s/%s is:issue", owner, repo)
	if q != "" {
		qString += fmt.Sprintf(" %s in:title", q)
	}
	qString = url.PathEscape(qString)
	u := fmt.Sprintf(static.GithubIssueSimpleSearchApi, qString, page, limit) // q page per_page
	log.Println(u,auth)
	bs, resp, err := g.httpService.HttpGet(u, func(r *http.Request) {
		r.Header.Set("Authorization", auth)
	})
	if err != nil {
		return 0, nil, err
	}
	log.Println(string(bs))
	if resp.StatusCode != 200 {
		return 0, nil, errors.New("response status is not 200 OK")
	}

	r := &biz.GithubIssueSearchResult{}
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

func (g *GithubService) GetIssueTimelines(name string, page int32, auth string) ([]map[string]interface{}, error) {
	// get user related issues
	issueUrls := make([]string, 0)
	issueCts := make([]string, 0)
	issueUsers := make([]map[string]interface{}, 0)
	getIssues := func(page int) (urls []string, cts []string, users []map[string]interface{}, tot int32, err error) {
		url := fmt.Sprintf(static.GithubIssueSearchApi, "updated", "desc", name, page, 100) // sort order involve page per_page
		bs, _, err := g.httpService.HttpGet(url, func(r *http.Request) {
			r.Header.Add("Authorization", auth)
			r.Header.Add("Accept", static.GithubAcceptPreview)
		})
		if err != nil {
			return nil, nil, nil, 0, err
		}

		data := &struct {
			TotalCount int32 `json:"total_count"`
			Items      []*struct {
				HtmlUrl   string                 `json:"html_url"`
				CreatedAt string                 `json:"created_at"`
				User      map[string]interface{} `json:"user"`
			} `json:"items"` // 30
		}{}
		err = json.Unmarshal(bs, data)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		urls = make([]string, len(data.Items))
		cts = make([]string, len(data.Items))
		users = make([]map[string]interface{}, len(data.Items))
		for idx := range data.Items {
			urls[idx] = data.Items[idx].HtmlUrl
			cts[idx] = data.Items[idx].CreatedAt
			users[idx] = data.Items[idx].User
		}

		return urls, cts, users, data.TotalCount, nil
	}

	pageUrls, pageCts, pageUsers, tot, err := getIssues(1)
	if err != nil {
		return nil, err
	}
	issueUrls = append(issueUrls, pageUrls...)
	issueCts = append(issueCts, pageCts...)
	issueUsers = append(issueUsers, pageUsers...)

	perPage := int32(len(issueUrls))
	pageCnt := int(math.Ceil(float64(tot) / float64(perPage)))
	enoughCnt := (page + 1) * static.GithubDefaultIssueLimit
	if perPage < enoughCnt && pageCnt > 1 { // not enough && has next page
		wg := sync.WaitGroup{}
		mu := sync.Mutex{}
		once := sync.Once{}
		errOnce := error(nil)
		for i := 2; i <= pageCnt; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				mu.Lock()
				l := len(issueUrls)
				mu.Unlock()
				if int32(l) >= enoughCnt { // enough
					return
				}
				pageUrls, pageCts, pageUsers, _, err := getIssues(i)
				if err != nil {
					once.Do(func() { errOnce = err })
					return
				}

				mu.Lock()
				issueUrls = append(issueUrls, pageUrls...)
				issueCts = append(issueCts, pageCts...)
				issueUsers = append(issueUsers, pageUsers...)
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		if errOnce != nil {
			return nil, errOnce
		}
	}
	if int32(len(issueUrls)) >= enoughCnt { // enough new
		issueUrls = issueUrls[:enoughCnt]
	}

	// parse issue url list
	type Issue struct {
		Owner     string
		Repo      string
		Number    int32
		CreatedAt string
		User      map[string]interface{}
	}
	issues := make([]*Issue, 0)
	for idx, url := range issueUrls {
		sp := strings.Split(url, "/") // https://github.com/gofiber/fiber/issues/556
		number, err := strconv.Atoi(sp[len(sp)-1])
		if err != nil {
			return nil, err
		}
		repo := sp[len(sp)-3]
		owner := sp[len(sp)-4]
		issues = append(issues, &Issue{
			Owner:     owner,
			Repo:      repo,
			Number:    int32(number),
			CreatedAt: issueCts[idx],
			User:      issueUsers[idx],
		})
	}

	// get issue events
	getIssuesCnt := 0
	getIssueTimeline := func(issue *Issue) ([]map[string]interface{}, error) {
		url := fmt.Sprintf(static.GithubIssueTimelineApi, issue.Owner, issue.Repo, issue.Number, 100)
		bs, _, err := g.httpService.HttpGet(url, func(r *http.Request) {
			r.Header.Add("Authorization", auth)
			r.Header.Add("Accept", static.GithubAcceptPreview)
		})
		if err != nil {
			return nil, err
		}
		getIssuesCnt++

		data := make([]map[string]interface{}, 0)
		err = json.Unmarshal(bs, &data)
		if err != nil {
			return nil, err
		}

		out := make([]map[string]interface{}, 0)
		for idx := range data {
			if data[idx]["event"] == "subscribed" { // filter subscribed, preserve mentioned
				continue
			}
			data[idx]["repo"] = fmt.Sprintf("%s/%s", issue.Owner, issue.Repo)
			data[idx]["number"] = issue.Number
			data[idx]["involve"] = name
			out = append(out, data[idx])
		}
		return out, nil
	}

	events := make([]map[string]interface{}, 0)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	once := sync.Once{}
	errOnce := error(nil)
	for _, issue := range issues {
		wg.Add(1)
		go func(issue *Issue, page int32) {
			defer wg.Done()

			data, e := getIssueTimeline(issue)
			if e != nil {
				once.Do(func() { errOnce = err })
				return
			}

			// append issue opened
			data = append([]map[string]interface{}{{
				"id":         nil, // <<< ATTENTION NIL
				"node_id":    nil,
				"event":      "opened",
				"actor":      issue.User,
				"commit_id":  nil,
				"commit_url": nil,
				"created_at": issue.CreatedAt,
				"repo":       fmt.Sprintf("%s/%s", issue.Owner, issue.Repo),
				"number":     issue.Number,
				"involve":    name,
				"url":        nil,
			}}, data...)

			// reverse
			for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
				data[i], data[j] = data[j], data[i]
			}

			mu.Lock()
			events = append(events, data...)
			mu.Unlock()
		}(issue, page)
	}
	wg.Wait()
	if errOnce != nil {
		return nil, errOnce
	}

	// filter issue event
	tempEvents := make([]map[string]interface{}, 0)
	for idx := range events {
		if _, ok := events[idx]["created_at"]; ok {
			tempEvents = append(tempEvents, events[idx])
		}
	}
	events = tempEvents

	sort.SliceStable(events, func(i, j int) bool {
		// create_at
		cti, oki := events[i]["created_at"]
		ctj, okj := events[j]["created_at"]
		if !oki {
			return false // j > i(x)
		}
		if !okj {
			return true // i > j(x)
		}
		ti, eri := time.Parse(time.RFC3339, cti.(string))
		tj, erj := time.Parse(time.RFC3339, ctj.(string))
		if eri != nil {
			return false // j > i(x)
		}
		if erj != nil {
			return true // i > j(x)
		}

		return ti.Unix() > tj.Unix()
	})

	l := int32(len(events))
	from := static.GithubDefaultIssueLimit * (page - 1)
	to := static.GithubDefaultIssueLimit * page
	if from >= l {
		return []map[string]interface{}{}, nil
	}
	if to > l {
		to = l
	}
	events = events[from:to]

	return events, nil
}
