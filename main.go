package main

import "log"

func main() {
	crawl := Crawler{URL: "https://parserdigital.com"}
	if err := crawl.Run(); err != nil {
		log.Fatal(err)
	}
}
