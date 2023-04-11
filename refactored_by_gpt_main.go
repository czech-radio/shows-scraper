package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

type program struct {
	name   string
	epName string
	time   string
}

func main() {
	// Define the desired time period
	startDate := time.Date(2021, 9, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2021, 9, 30, 23, 59, 59, 0, time.Local)

	// Create a map to store the programs
	programs := make(map[string][]program)

	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.irozhlas.cz"),
		colly.CacheDir("./cache"),
	)

	// Set up error logging
	logFile, err := os.Create("errors.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Find all program links on the given page
	c.OnHTML(".program-items__link", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Visit the program page
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Find the details of the programs on the given page
	c.OnHTML(".broadcasts__list__item", func(e *colly.HTMLElement) {
		// Extract the program name and episode name
		programName := e.ChildText(".broadcasts__list__item__header__title")
		epName := e.ChildText(".broadcasts__list__item__header__sub-title")

		// Extract the broadcast time
		broadcastTime, err := time.Parse("2.1.2006 15:04", e.ChildText(".broadcasts__list__item__header__time"))
		if err != nil {
			log.Println("Error parsing time:", err)
			return
		}

		// Check if the broadcast time is within the desired time period
		if broadcastTime.Before(startDate) || broadcastTime.After(endDate) {
			return
		}

		// Add the program to the map
		p := program{name: programName, epName: epName, time: broadcastTime.Format("2006-01-02 15:04:05")}
		programs[programName] = append(programs[programName], p)

		// Print the program details
		fmt.Printf("Program: %s\nEpisode: %s\nTime: %s\n\n", programName, epName, broadcastTime.Format("2006-01-02 15:04:05"))
	})

	// Visit the main page of the publicistika section
	c.Visit("https://www.irozhlas.cz/publicistika")

	// Print the programs and their broadcasts
	for programName, programList := range programs {
		fmt.Printf("Program: %s\n", programName)
		for _, p := range programList {
			fmt.Printf("Episode: %s\nTime: %s\n\n", p.epName, p.time)
		}
	}
}
