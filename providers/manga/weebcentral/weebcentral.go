package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tsuki-reader/nisshoku/providers"
)

type WeebcentralProvider struct {
	context providers.ProviderContext
	client  *http.Client
	baseUrl string
}

func NewProvider(context providers.ProviderContext) providers.Provider {
	var provider WeebcentralProvider
	provider.context = context
	provider.client = &http.Client{
		Timeout: 30 * time.Second,
	}
	provider.baseUrl = "https://weebcentral.com"
	return &provider
}

func (p *WeebcentralProvider) Search(query string) ([]providers.ProviderResult, error) {
	resp, err := p.retrieveHtml("/search/data?author=&text=" + query + "&sort=Best%20Match&order=Ascending&official=Any&display_mode=Full%20Display")
	if err != nil {
		return []providers.ProviderResult{}, err
	}
	defer resp.Body.Close()

	doc, err := p.parseDocument(resp)
	if err != nil {
		return []providers.ProviderResult{}, err
	}

	articles := doc.Find("article.bg-base-300")
	results := []providers.ProviderResult{}
	for i := range articles.Nodes {
		article := articles.Eq(i)
		header := article.Find("section.hidden > div.text-lg > span > a")
		title := header.Text()
		resultUrl, _ := header.Attr("href")
		id, err := p.extractId(resultUrl)
		if err != nil {
			// Just eff it off
			continue
		}
		provider := "weebcentral"
		startYear, err := strconv.Atoi(article.Find("section.hidden > div:nth-child(2) > span").Text())
		if err != nil {
			startYear = 0
		}

		result := providers.ProviderResult{
			Title:            title,
			ID:               id,
			Provider:         provider,
			AlternativeNames: []string{},
			StartYear:        startYear,
		}
		results = append(results, result)
	}

	return results, nil
}

func (p *WeebcentralProvider) GetChapters(id string) ([]providers.Chapter, error) {
	fullChapterListPath := getFullChapterListUri(id)
	resp, err := p.retrieveHtml(fullChapterListPath)
	if err != nil {
		return []providers.Chapter{}, err
	}
	defer resp.Body.Close()

	doc, err := p.parseDocument(resp)
	if err != nil {
		return []providers.Chapter{}, err
	}

	anchors := doc.Find("body > a.flex")
	count := len(anchors.Nodes)
	results := []providers.Chapter{}
	for i := range anchors.Nodes {
		anchor := anchors.Eq(i)
		absoluteNumber := count - i
		title := anchor.Find("span.grow.flex.items-center.gap-2 > span:nth-child(1)").Text()
		chapterUri, _ := anchor.Attr("href")
		id, err := p.extractId(chapterUri)
		if err != nil {
			// This should never happen. If it does then something has gone HORRIBLY wrong.
			continue
		}

		result := providers.Chapter{
			Title:          title,
			ID:             id,
			Provider:       "weebcentral",
			AbsoluteNumber: absoluteNumber,
		}
		results = append(results, result)
	}

	return results, nil
}

func (p *WeebcentralProvider) GetChapterPages(id string) ([]providers.Page, error) {
	resp, err := p.retrieveHtml(id + "/images?is_prev=False&current_page=1&reading_style=long_strip")
	if err != nil {
		return []providers.Page{}, err
	}
	defer resp.Body.Close()

	doc, err := p.parseDocument(resp)
	if err != nil {
		return []providers.Page{}, err
	}

	images := doc.Find("section.flex > img")
	results := []providers.Page{}
	for i := range images.Nodes {
		image := images.Eq(i)
		src, _ := image.Attr("src")

		page := providers.Page{
			Provider:   "weebcentral",
			ImageURL:   src,
			PageNumber: i,
		}
		results = append(results, page)
	}

	return results, nil
}

func (p *WeebcentralProvider) ImageHeaders() map[string]string {
	return map[string]string{}
}

func (p *WeebcentralProvider) ProviderType() providers.ProviderType {
	return providers.Manga
}

func (p *WeebcentralProvider) retrieveHtml(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", p.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (p *WeebcentralProvider) parseDocument(response *http.Response) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(response.Body)
}

func (p *WeebcentralProvider) extractId(uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return parsed.Path, nil
}

func getFullChapterListUri(id string) string {
	lastSlashIndex := strings.LastIndex(id, "/")
	if lastSlashIndex != -1 {
		return id[:lastSlashIndex] + "/full-chapter-list"
	} else {
		// TODO: Don't silent fail
		return id + "/full-chapter-list"
	}
}
