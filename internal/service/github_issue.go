package service

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/common_api/internal/pkg/static"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (g *GithubService) GetIssueTimelines(name string, page uint32, auth string) ([]map[string]any, error) {
	// get user related issues
	issueUrls := make([]string, 0)
	issueCts := make([]string, 0)
	issueUsers := make([]map[string]any, 0)
	getIssues := func(page int) (urls []string, cts []string, users []map[string]any, tot int32, err error) {
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
				HtmlUrl   string         `json:"html_url"`
				CreatedAt string         `json:"created_at"`
				User      map[string]any `json:"user"`
			} `json:"items"` // 30
		}{}
		err = json.Unmarshal(bs, data)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		urls = make([]string, len(data.Items))
		cts = make([]string, len(data.Items))
		users = make([]map[string]any, len(data.Items))
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

	perPage := uint32(len(issueUrls))
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
				if uint32(l) >= enoughCnt { // enough
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
	if uint32(len(issueUrls)) >= enoughCnt { // enough new
		issueUrls = issueUrls[:enoughCnt]
	}

	// parse issue url list
	type Issue struct {
		Owner     string
		Repo      string
		Number    int32
		CreatedAt string
		User      map[string]any
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
	getIssueTimeline := func(issue *Issue) ([]map[string]any, error) {
		url := fmt.Sprintf(static.GithubIssueTimelineApi, issue.Owner, issue.Repo, issue.Number, 100)
		bs, _, err := g.httpService.HttpGet(url, func(r *http.Request) {
			r.Header.Add("Authorization", auth)
			r.Header.Add("Accept", static.GithubAcceptPreview)
		})
		if err != nil {
			return nil, err
		}
		getIssuesCnt++

		data := make([]map[string]any, 0)
		err = json.Unmarshal(bs, &data)
		if err != nil {
			return nil, err
		}

		out := make([]map[string]any, 0)
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

	events := make([]map[string]any, 0)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	once := sync.Once{}
	errOnce := error(nil)
	for _, issue := range issues {
		wg.Add(1)
		go func(issue *Issue, page uint32) {
			defer wg.Done()

			data, e := getIssueTimeline(issue)
			if e != nil {
				once.Do(func() { errOnce = err })
				return
			}

			// append issue opened
			data = append([]map[string]any{{
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
	tempEvents := make([]map[string]any, 0)
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

	l := uint32(len(events))
	from := static.GithubDefaultIssueLimit * (page - 1)
	to := static.GithubDefaultIssueLimit * page
	if from >= l {
		return []map[string]any{}, nil
	}
	if to > l {
		to = l
	}
	events = events[from:to]

	return events, nil
}
