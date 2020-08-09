package service

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/common/logger"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"math"
	"net/http"
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
		httpService: xdi.GetByNameForce(sn.SHttpService).(*HttpService),
	}
}

// noinspection GoSnakeCaseUsage
const (
	GITHUB_RATE_LIMIT_URL     = "https://api.github.com/rate_limit"
	GITHUB_SEARCH_ISSUE_URL   = "https://api.github.com/search/issues?sort=%s&order=%s&q=involves:%s&page=%d&per_page=%d"
	GITHUB_ISSUE_TIMELINE_URL = "https://api.github.com/repos/%s/%s/issues/%d/timeline?per_page=%d"
	GITHUB_ISSUE_EVENT_LIMIT  = 20
)

func (g *GithubService) GetRateLimit(auth string) (map[string]interface{}, error) {
	header := &http.Header{}
	header.Add("Authorization", auth)
	header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpService.HttpGet(GITHUB_RATE_LIMIT_URL, header, nil)
	if err != nil {
		return nil, err
	}

	core := make(map[string]interface{})
	err = json.Unmarshal(resp, &core)
	if err != nil {
		return nil, err
	}
	return core, nil
}

func (g *GithubService) GetIssueEvents(name string, page int32, auth string) ([]map[string]interface{}, error) {
	header := &http.Header{}
	header.Add("Authorization", auth)
	header.Add("Accept", "application/vnd.github.mockingbird-preview+json")

	// get user related issues
	issueUrls := make([]string, 0)
	issueCts := make([]string, 0)
	issueUsers := make([]map[string]interface{}, 0)
	getIssues := func(page int) (urls []string, cts []string, users []map[string]interface{}, tot int32, err error) {
		url := fmt.Sprintf(GITHUB_SEARCH_ISSUE_URL, "updated", "desc", name, page, 100) // sort order involve page per_page
		bs, err := g.httpService.HttpGet(url, header, logger.LogGhUrl)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		data := &struct {
			TotalCount int32 `json:"total_count"`
			Items      []*struct { // 30
				HtmlUrl   string                 `json:"html_url"`
				CreatedAt string                 `json:"created_at"`
				User      map[string]interface{} `json:"user"`
			} `json:"items"`
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
	enoughCnt := (page + 1) * GITHUB_ISSUE_EVENT_LIMIT
	if perPage < enoughCnt && pageCnt > 1 { // not enough && has next page
		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}
		wg.Add(pageCnt - 1)

		var err error
		for i := 2; i <= pageCnt; i++ {
			go func(i int, err *error) {
				if *err != nil || int32(len(issueUrls)) >= enoughCnt { // enough
					wg.Done()
					return
				}

				pageUrls, pageCts, pageUsers, _, e := getIssues(i)
				if e != nil {
					mu.Lock()
					*err = e
					mu.Unlock()
					wg.Done()
					return
				}
				mu.Lock()
				issueUrls = append(issueUrls, pageUrls...)
				issueCts = append(issueCts, pageCts...)
				issueUsers = append(issueUsers, pageUsers...)
				mu.Unlock()
				wg.Done()
			}(i, &err)
		}
		wg.Wait()
		if err != nil {
			return nil, err
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
		url := fmt.Sprintf(GITHUB_ISSUE_TIMELINE_URL, issue.Owner, issue.Repo, issue.Number, 100)
		bs, err := g.httpService.HttpGet(url, header, nil)
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
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(issues))

	err = nil
	for _, issue := range issues {
		go func(issue *Issue, page int32, err *error) {
			if *err != nil {
				wg.Done()
				return
			}

			data, e := getIssueTimeline(issue)
			if e != nil {
				mu.Lock()
				*err = e
				mu.Unlock()
				wg.Done()
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
			wg.Done()
		}(issue, page, &err)
	}
	wg.Wait()
	if err != nil {
		return nil, err
	}
	logger.LogGhUrl(fmt.Sprintf("get issue event count: %d", getIssuesCnt))

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
	from := GITHUB_ISSUE_EVENT_LIMIT * (page - 1)
	to := GITHUB_ISSUE_EVENT_LIMIT * page
	if from >= l {
		return []map[string]interface{}{}, nil
	}
	if to > l {
		to = l
	}
	events = events[from:to]

	return events, nil
}
