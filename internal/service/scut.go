package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"github.com/Aoi-hosizora/ahlib/xerror"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/model/obj"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/static"
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
		httpService: xmodule.MustGetByName(sn.SHttpService).(*HttpService),
	}
}

func (s *ScutService) GetJwNotices() ([]*obj.ScutNoticeItem, error) {
	form := &url.Values{}
	form.Add("category", "0")
	form.Add("tag", "0") // all tags
	form.Add("pageNum", "1")
	form.Add("pageSize", "50")
	form.Add("keyword", "")
	bs, _, err := s.httpService.HttpPost(static.ScutJwApi, strings.NewReader(form.Encode()), func(r *http.Request) {
		r.Header.Set(headers.Cookie, static.ScutJwCookie)
		r.Header.Set(headers.ContentType, static.ScutJwContentType)
		r.Header.Set(headers.Origin, static.ScutJwOrigin)
		r.Header.Set(headers.Referer, static.ScutJwReferer)
	})
	if err != nil {
		return nil, err
	}

	type Result struct {
		List []*struct {
			Id         string `json:"id"`
			CreateTime string `json:"createTime"`
			Title      string `json:"title"`
			IsNew      bool   `json:"isNew"`
			Tag        int    `json:"tag"`
		} `json:"list"`
	}
	items := &Result{}
	err = json.Unmarshal(bs, items)
	if err != nil {
		return nil, err
	}

	out := make([]*obj.ScutNoticeItem, 0, len(items.List))
	for _, item := range items.List {
		u := fmt.Sprintf(static.ScutJwNoticeUrl, item.Id)
		mu := fmt.Sprintf(static.ScutJwNoticeMobileUrl, item.Id)
		t := strings.ReplaceAll(item.CreateTime, ".", "-") // 2020.01.01 -> 2020-01-01
		out = append(out, obj.NewScutNoticeItem(item.Title, u, mu, static.ScutJwTagNames[item.Tag], t))
	}
	obj.SortScutNoticeItems(out)
	return out, nil
}

type scutTagWithDoc struct {
	tag string
	doc *goquery.Document
}

func (s *ScutService) GetSeNotices() ([]*obj.ScutNoticeItem, error) {
	pairs := make([]scutTagWithDoc, len(static.ScutSeTagParts))
	eg := xerror.NewErrorGroup(context.Background())
	for idx, part := range static.ScutSeTagParts {
		pairs[idx].tag = static.ScutSeTagNames[idx]
		idx, part := idx, part
		eg.Go(func(ctx context.Context) error {
			bs, _, err := s.httpService.HttpGetWithCtx(ctx, fmt.Sprintf(static.ScutSeWebUrl, part), nil)
			if err != nil {
				return err
			}
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
			if err != nil {
				return err
			}
			pairs[idx].doc = doc
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	out := make([]*obj.ScutNoticeItem, 0)
	for _, pair := range pairs {
		pair.doc.Find("ul.news_ul > li.news_li").Each(func(i int, li *goquery.Selection) {
			a := li.Find("span.news_title a")
			meta := li.Find("span.news_meta") // 2019-10-01
			href := fmt.Sprintf(static.ScutSeNoticeWebUrl, a.AttrOr("href", ""))
			out = append(out, obj.NewScutNoticeItem(a.Text(), href, href, "软院"+pair.tag, meta.Text()))
		})
	}
	obj.SortScutNoticeItems(out)
	return out, nil
}

func (s *ScutService) GetGrNotices() ([]*obj.ScutNoticeItem, error) {
	const pages = 2
	docs := make([]*goquery.Document, pages)
	eg := xerror.NewErrorGroup(context.Background())
	for i := 0; i < pages; i++ {
		i := i
		eg.Go(func(ctx context.Context) error {
			bs, _, err := s.httpService.HttpGetWithCtx(ctx, fmt.Sprintf(static.ScutGrWebUrl, i+1), nil)
			if err != nil {
				return err
			}
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
			if err != nil {
				return err
			}
			docs[i] = doc
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	out := make([]*obj.ScutNoticeItem, 0)
	for _, doc := range docs {
		doc.Find("table.wp_article_list_table tr").Each(func(i int, tr *goquery.Selection) {
			a := tr.Find("a")
			href := fmt.Sprintf(static.ScutGrNoticeWebUrl, a.AttrOr("href", ""))
			span := a.Find("span")
			t := span.Text() // 2022-01-19
			span.Remove()
			out = append(out, obj.NewScutNoticeItem(a.Text(), href, href, "研究生院"+static.ScutGrTagName, t))
		})
	}
	obj.SortScutNoticeItems(out)
	return out, nil
}

func (s *ScutService) GetGzicNotices() ([]*obj.ScutNoticeItem, error) {
	pairs := make([]scutTagWithDoc, len(static.ScutGzicTagParts))
	eg := xerror.NewErrorGroup(context.Background())
	for idx, part := range static.ScutGzicTagParts {
		pairs[idx].tag = static.ScutGzicTagNames[idx]
		idx, part := idx, part
		eg.Go(func(ctx context.Context) error {
			bs, _, err := s.httpService.HttpGetWithCtx(ctx, fmt.Sprintf(static.ScutGzicWebUrl, part), nil)
			if err != nil {
				return err
			}
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
			if err != nil {
				return err
			}
			pairs[idx].doc = doc
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	out := make([]*obj.ScutNoticeItem, 0)
	for _, pair := range pairs {
		pair.doc.Find("div.right-nr div.row div.thr-box").Each(func(i int, div *goquery.Selection) {
			a := div.Find("a")
			href := a.AttrOr("href", "")
			if strings.HasPrefix(href, "/gzic") {
				href = fmt.Sprintf(static.ScutGzicNoticeWebUrl, href)
			}
			span := a.Find("span") // 2022-01-17
			p := a.Find("p")
			out = append(out, obj.NewScutNoticeItem(p.Text(), href, href, "GZIC"+pair.tag, span.Text()))
		})
	}
	obj.SortScutNoticeItems(out)
	return out, nil
}
