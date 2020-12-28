package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/model/vo"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/common_api/src/static"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

type ScutService struct {
	httpService *HttpService
}

func NewScutService() *ScutService {
	return &ScutService{
		httpService: xdi.GetByNameForce(sn.SHttpService).(*HttpService),
	}
}

func (s *ScutService) GetJwItems() ([]*vo.ScutPostItem, error) {
	form := &url.Values{}
	form.Add("tag", "0")
	form.Add("pageNum", "1")
	form.Add("pageSize", "50")
	form.Add("keyword", "")
	bs, _, err := s.httpService.HttpPost(static.SCUT_JW_API_URL, strings.NewReader(form.Encode()), func(r *http.Request) {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("User-Agent", static.SCUT_JW_USER_AGENT)
		r.Header.Set("Referer", static.SCUT_JW_REFERER)
	})
	if err != nil {
		return nil, err
	}

	items := &struct {
		List []*struct {
			CreateTime string `json:"createTime"`
			Id         string `json:"id"`
			IsNew      bool   `json:"isNew"`
			Tag        int    `json:"tag"`
			Title      string `json:"title"`
		} `json:"list"`
		Message string `json:"message"`
		PageNum int    `json:"pagenum"`
		Row     int    `json:"row"`
		Success bool   `json:"success"`
		Total   int    `json:"total"`
	}{}
	err = json.Unmarshal(bs, items)
	if err != nil {
		return nil, err
	}

	out := make([]*vo.ScutPostItem, len(items.List))
	for i, item := range items.List {
		out[i] = &vo.ScutPostItem{
			Title:     item.Title,
			Url:       fmt.Sprintf(static.SCUT_JW_ITEM_URL, item.Id),
			MobileUrl: fmt.Sprintf(static.SCUT_JW_ITEM_MOBILE_URL, item.Id),
			Type:      static.SCUT_JW_TAG_NAMES[item.Tag-1],
			Date:      strings.ReplaceAll(item.CreateTime, ".", "-"), // 2020-01-01
		}
	}

	return out, nil
}

func (s *ScutService) GetSeItems() ([]*vo.ScutPostItem, error) {
	type Item struct {
		TagName string            // <<
		Doc     *goquery.Document // <<
	}

	items := make([]*Item, len(static.SCUT_SE_TAG_PARTS))
	for idx, part := range static.SCUT_SE_TAG_PARTS {
		u := fmt.Sprintf(static.SCUT_SE_WEB_URL, part)
		bs, _, err := s.httpService.HttpGet(u, func(r *http.Request) {
			r.Header.Set("User-Agent", static.SCUT_SE_USER_AGENT)
		})
		if err != nil {
			return nil, err
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
		if err != nil {
			return nil, err
		}
		items[idx] = &Item{TagName: static.SCUT_SE_TAG_NAMES[idx], Doc: doc}
	}

	out := make([]*vo.ScutPostItem, 0)
	for _, item := range items {
		lis := item.Doc.Find("ul.news_ul > li.news_li")
		news := make([]*vo.ScutPostItem, lis.Size())
		lis.Each(func(i int, s *goquery.Selection) {
			a := s.Find(".news_title a")
			u := fmt.Sprintf(static.SCUT_SE_ITEM_URL, a.AttrOr("href", ""))
			meta := s.Find("span.news_meta")
			news[i] = &vo.ScutPostItem{
				Title:     a.Text(),
				Url:       u,
				MobileUrl: u,
				Type:      "软院" + item.TagName,
				Date:      meta.Text(), // 2019-10-01
			}
		})
		out = append(out, news...)
	}

	return out, nil
}
