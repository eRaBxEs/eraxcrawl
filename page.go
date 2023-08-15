package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL   string
	Title string
	Links []string

	response *http.Response
}

func (p *Page) Goto() (err error) {
	cli := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		return fmt.Errorf("page.Goto: failed to create request %w", err)
	}

	p.response, err = cli.Do(req)
	if err != nil {
		return fmt.Errorf("page.Goto: failed to get page %w", err)
	}

	return nil
}

func (p *Page) GetLinks() (err error) {
	doc, err := goquery.NewDocumentFromReader(p.response.Body)
	if err != nil {
		return fmt.Errorf("page.GetLinks: failed to parse page %w", err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link := s.AttrOr("href", "")
		if !p.linkInList(link) {
			p.Links = append(p.Links, link)
			fmt.Println(link)
		}

	})

	return nil
}

func (p *Page) linkInList(link string) bool {
	var found bool
	for _, l := range p.Links {
		found = strings.EqualFold(link, l)
		if found {
			return true
		}
	}

	return false
}
