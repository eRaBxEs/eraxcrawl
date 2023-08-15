package main

import (
	"fmt"
	"log"
)

type Crawler struct {
	URL   string
	pages []*Page
}

func (c *Crawler) Run() error {
	if err := c.Crawl(c.URL); err != nil {
		return err
	}

	c.Output()

	return nil
}

func (c *Crawler) Crawl(url string) (err error) {

	// channelLinks := make(chan []string, 5)
	page := &Page{URL: url}
	log.Println("Crawling:", url)

	err = page.Goto()
	if err != nil {
		return fmt.Errorf("page.GoTo: error occured: %w", err)
	}

	log.Println("Getting links from:", url)

	err = page.GetLinks()
	if err != nil {
		return fmt.Errorf("page.GetLinks: Error occured: %v", err)
	}

	c.pages = append(c.pages, page)

	if page.Links == nil {
		return nil
	}

	// fmt.Println("Crawl:", c.pages[len(c.pages)-1].Links)
	strLength := c.pages[len(c.pages)-1].Links[1:]

	for i := 0; i < len(strLength); i++ {
		go c.Output()
		c.Crawl(strLength[i])
	}

	return nil
}

func (c *Crawler) Output() {

	for _, p := range c.pages {
		fmt.Println(p.Title)
		for i, l := range p.Links {
			fmt.Println(i, ":", l)
		}

	}
}
