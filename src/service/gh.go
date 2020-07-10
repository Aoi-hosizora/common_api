package service

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/common_api/src/common/logger"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// noinspection GoSnakeCaseUsage
const (
	GITHUB_SEARCH_ISSUE_URL  = "https://api.github.com/search/issues?sort=created&order=desc&q=involves:%s&page=%d"
	GITHUB_ISSUE_EVENT_URL   = "https://api.github.com/repos/%s/%s/issues/%d/events?page=%d"
	GITHUB_ISSUE_EVENT_PAGES = 3
	GITHUB_ISSUE_EVENT_LIMIT = 30
)

func GetIssueEvents(name string, page int32, auth string) ([]map[string]interface{}, error) {
	header := &http.Header{}
	if auth != "" {
		header.Add("Authorization", auth)
	}

	// get user related issues
	urls := make([]string, 0)
	getIssues := func(page int) ([]string, int32, error) {
		url := fmt.Sprintf(GITHUB_SEARCH_ISSUE_URL, name, page)
		bs, err := HttpGet(url, header, logger.LogGhUrl)
		if err != nil {
			return nil, 0, err
		}

		data := &struct {
			TotalCount int32 `json:"total_count"`
			Items      []*struct { // 30
				HtmlUrl string `json:"html_url"`
			} `json:"items"`
		}{}
		err = json.Unmarshal(bs, data)
		if err != nil {
			return nil, 0, err
		}

		urls := make([]string, len(data.Items))
		for idx := range data.Items {
			urls[idx] = data.Items[idx].HtmlUrl
		}

		return urls, data.TotalCount, nil
	}

	pageUrls, tot, err := getIssues(1)
	if err != nil {
		return nil, err
	}
	urls = append(urls, pageUrls...)

	cnt := int(math.Ceil(float64(tot) / float64(len(pageUrls))))
	if cnt > 1 {
		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}
		wg.Add(cnt - 1)

		var err error
		for i := 2; i <= cnt; i++ {
			go func(i int, err *error) {
				if *err != nil {
					wg.Done()
					return
				}

				pageUrls, _, e := getIssues(i)
				if e != nil {
					mu.Lock()
					*err = e
					mu.Unlock()
					wg.Done()
					return
				}
				mu.Lock()
				urls = append(urls, pageUrls...)
				mu.Unlock()
				wg.Done()
			}(i, &err)
		}
		wg.Wait()
		if err != nil {
			return nil, err
		}
	}

	// parse issue url list
	type Issue struct {
		Owner  string
		Repo   string
		Number int32
	}

	issues := make([]*Issue, 0)
	for _, url := range urls {
		sp := strings.Split(url, "/") // https://github.com/gofiber/fiber/issues/556
		number, err := strconv.Atoi(sp[len(sp)-1])
		if err != nil {
			return nil, err
		}
		repo := sp[len(sp)-3]
		owner := sp[len(sp)-4]
		issues = append(issues, &Issue{
			Owner:  owner,
			Repo:   repo,
			Number: int32(number),
		})
	}

	// get issue events
	getIssuesCnt := 0
	getIssueEvents := func(issue *Issue, page int32) ([]map[string]interface{}, error) {
		url := fmt.Sprintf(GITHUB_ISSUE_EVENT_URL, issue.Owner, issue.Repo, issue.Number, page)
		// bs, err := HttpGet(url, header, logger.LogGhUrl)
		bs, err := HttpGet(url, header, nil)
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
			if data[idx]["event"] == "subscribed" {
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
	wg.Add(len(issues) * GITHUB_ISSUE_EVENT_PAGES)

	err = nil
	for _, issue := range issues {
		emptied := false
		for page := int32(1); page <= GITHUB_ISSUE_EVENT_PAGES; page++ {
			go func(issue *Issue, page int32, err *error, emptied *bool) {
				if *err != nil || *emptied {
					wg.Done()
					return
				}

				data, e := getIssueEvents(issue, page)
				if e != nil {
					mu.Lock()
					*err = e
					mu.Unlock()
					wg.Done()
					return
				} else if len(data) == 0 {
					*emptied = true
				}

				mu.Lock()
				events = append(events, data...)
				mu.Unlock()
				wg.Done()
			}(issue, page, &err, &emptied)
		}
	}
	wg.Wait()
	if err != nil {
		return nil, err
	}
	logger.LogGhUrl(fmt.Sprintf("get issue event count: %d", getIssuesCnt))

	// filter issue event
	sort.Slice(events, func(i, j int) bool {
		cti, oki := events[i]["created_at"]
		ctj, okj := events[j]["created_at"]
		if !oki || !okj {
			return false
		}
		ti, eri := time.Parse(time.RFC3339, cti.(string))
		tj, erj := time.Parse(time.RFC3339, ctj.(string))
		return eri == nil && erj == nil && ti.Unix() > tj.Unix()
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
