package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Crawl(t *testing.T) {
	testCase := []struct {
		name     string
		url      string
		response *http.Response
		link     []string
	}{
		{
			name:     "start crawl",
			url:      "http://start.com",
			response: &http.Response{Body: NewMockReaderCloser(testResponses[0])},
			link: []string{
				"http://siteA.com",
				"http://siteB.com",
				"http://siteC.com",
			},
		},
	}
	expectedCase := []struct {
		url   string
		pages []Page
	}{
		{
			url: "http://start.com",
			pages: []Page{
				{
					URL: "http://start.com",
					Links: []string{
						"http://siteA.com",
						"http://siteB.com",
						"http://siteC.com",
					},
				},
			},
		},
	}

	// Test to satisfy the GoTo called in Crawl
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	page := Page{URL: srv.URL}
	if err := page.Goto(); err != nil {
		t.Fatal(err)
	}

	if page.response.StatusCode != http.StatusOK {
		t.Errorf("page: %s returned %d", page.URL, page.response.StatusCode)
	}
	var crawledPage []Page
	// Test to satisfy the GetLinks called in Crawl
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			page := Page{response: tc.response}
			if err := page.GetLinks(); err != nil {
				t.Fatal(err)
			}

			if len(page.Links) != len(tc.link) {
				t.Error("incorrect number of links - expected:", len(tc.link), " got:", len(page.Links))
			}
			for i, l := range page.Links {
				if tc.link[i] != l {
					t.Error("incorrect links - expected:", tc.link[i], " got:", l)

				}
			}
			crawledPage = append(crawledPage, page)
		})
	}

	if len(expectedCase[0].pages) != len(crawledPage) {
		t.Error("incorrect number of pages - expected:", len(expectedCase[0].pages), "got:", len(crawledPage))
	}

}
