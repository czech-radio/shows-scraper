package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	links := []string{
		"https://www.radio.cz/cz/rubrika/zpravy",
		"https://www.radio.cz/cz/rubrika/udalosti",
		"https://www.radio.cz/cz/rubrika/zahranici",
	}

	for _, link := range links {
		fmt.Printf("Downloading articles from %s\n", link)
		resp, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Find all article links on the page and download them
		doc.Find(".list-news a").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if !ok {
				return
			}

			fmt.Printf("Downloading article from %s\n", href)
			resp, err := http.Get(href)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			// TODO: process the article content here
		})
	}

	fmt.Println("All articles downloaded successfully")
}
