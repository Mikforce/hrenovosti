package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page.
	res, err := http.Get("https://panorama.pub")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".grid-cols-1 ").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		fmt.Fprintf(w, "Review %d: %s\n", i, title)
	})
}

func main() {
	http.HandleFunc("/", ExampleScrape)
	log.Fatal(http.ListenAndServe(":80", nil))
}
