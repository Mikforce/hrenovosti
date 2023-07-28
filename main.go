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

	// Create an HTML string to store the parsed data
	html := "<html><body>"

	// Find the news items
	doc.Find(".grid-cols-1").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title and image URL
		title := s.Find("a").Text()
		imageURL, _ := s.Find("img").Attr("src")

		// Append the title and image URL to the HTML string
		html += fmt.Sprintf("<h2>%s</h2>", title)
		html += fmt.Sprintf(`<img src="%s" alt="%s">`, imageURL, title)
	})

	html += "</body></html>"

	// Set the response headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Write the HTML string to the response
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/", ExampleScrape)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
