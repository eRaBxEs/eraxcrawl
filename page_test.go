package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPage_Goto(t *testing.T) {

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

}

func TestPage_GetLinks(t *testing.T) {
	testCase := []struct {
		name     string
		response *http.Response
		link     []string
	}{
		{
			name:     "get links",
			response: &http.Response{Body: NewMockReaderCloser(testResponses[0])},
			link: []string{
				"http://siteA.com",
				"http://siteB.com",
				"http://siteC.com",
			},
		},
	}

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
		})
	}

}

var testResponses = []string{
	`
 <div>
	<h1>hello World</h1>
	<a href="http://siteA.com">Link 1</a>
	<a href="http://siteB.com">Link 2</a>
	<a href="http://siteC.com">Link 3</a>	
	<a href="http://siteC.com">Link 4</a>	
 </div>
 `,
}

// for the purpose of mocking
type mockReaderCloser struct {
	strings.Reader
}

func NewMockReaderCloser(s string) *mockReaderCloser {
	mc := new(mockReaderCloser)
	mc.Reader = *strings.NewReader(s)
	return mc
}

func (mockReaderCloser) Close() error {
	return nil
}
