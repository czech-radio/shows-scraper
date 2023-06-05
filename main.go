package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

var (
	versionFlag     bool
	buildTime       string
	sha1GitRevision string
	versionGitTag   string
)

var logger = log.New(os.Stderr, "", 0)

func main() {

	numPages := flag.Int("p", 1, "Number of pages to download.")
	flag.BoolVar(&versionFlag, "v", false, "Print application version and exit.")
	flag.Parse()

	// Must be at top!
	if versionFlag {
		fmt.Printf("Version: %s %s %s\n", versionGitTag, buildTime, sha1GitRevision)
		os.Exit(0)
	}

	for i := 0; i < *numPages; i++ {
		logger.Println("Page: ", fmt.Sprintf("%d/%d", i+1, *numPages))
		logger.Println("----------")

		episodes := GetRozhovoryEpisodes(i)
		sortArticlesByDate(episodes)

		for _, episode := range episodes {
			logger.Println(episode.Date)
			logger.Println("----------")
			episodeJSON, err := json.Marshal(episode)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(episodeJSON))
		}
	}
}

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Guest of some episode.
type Guest struct {
	Person
	Function string `json:"function"`
}

// Article scraped from website.
type Article struct {
	Show        string  `json:"show"`
	Episode     string  `json:"episode"`
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Link        string  `json:"link"`
	Teaser      string  `json:"teaser"`
	Moderator   Person  `json:"moderator"`
	Guests      []Guest `json:"guests"`
	TopicsCount int     `json:"topicsCount"`
}

// NewArticle creates a new article.
func NewArticle(show string, episode string, date string, time string, link string, teaser string, moderator Person, guests []Guest, topicsCount int) Article {
	article := Article{}
	article.Show = show
	article.Episode = episode
	article.Date = date
	article.Time = time
	article.Link = link
	article.Teaser = teaser
	article.Moderator = moderator
	article.Guests = guests
	article.TopicsCount = topicsCount

	return article
}

// sortArticlesByDate sorts articles in-place by it's date.
func sortArticlesByDate(articles []Article) {
	sort.SliceStable(articles, func(i, j int) bool {
		ci, cj := fmt.Sprintf("%s %s", articles[i].Show, articles[i].Date), fmt.Sprintf("%s %s", articles[j].Show, articles[j].Date)

		switch {
		case ci != cj:
			return ci > cj
		default:
			return ci > cj
		}
	})
}

// convertDate converts date string from e.g. `18. duben 2023`to `2023-04-18` format.
func convertDate(input string) string {
	s := strings.Split(input, " ")
	day, err := strconv.Atoi(strings.Split(s[0], ".")[0])
	if err != nil {
		// log.Println(fmt.Sprintf("Couldn't get day from date: %s", err.Error()))
		return input

	}
	year := s[2]
	months := map[string]string{
		"leden":    "01",
		"únor":     "02",
		"březen":   "03",
		"duben":    "04",
		"květen":   "05",
		"červen":   "06",
		"červenec": "07",
		"srpen":    "08",
		"září":     "09",
		"říjen":    "10",
		"listopad": "11",
		"prosinec": "12",
	}
	mo := months[s[1]]

	return fmt.Sprintf("%s-%s-%02d", year, mo, day)
}

// normalizeNames converts given personal name to normalized form.
// We use hard coded moderator lookup names table e.g `Tomáše Pancíře` -> `Tomáš Pancíř`.
func normalizeNames(name string) string {

	name = strings.TrimSpace(name)

	if name == "Tomáše Pancíře" {
		name = "Tomáš Pancíř"
	}
	if name == "Věry Štechrová" {
		name = "Věra Štechrová"
	}
	if name == "Vladimíra Kroce" {
		name = "Vladimír Kroc"
	}
	if name == "Tomáše Pavlíčka" {
		name = "Tomáš Pavlíček"
	}

	if name == "Petra Dudka" {
		name = "Petr Dudek"
	}

	return name
}

func countTopics(title string) int {
	return len(strings.Split(title, "."))
}

// GetRozhovoryEpisodes scrapes *Hlavní zprávy - rozhovory, komentáře* episodes data from
// https://radiozurnal.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846.
func GetRozhovoryEpisodes(pageNumber int) []Article {

	show := "Hlavní zprávy - rozhovory a komentáře"
	var teaser, episode, date, time string
	var moderator Person
	var moderator2 string
	var guests []Guest
	var topicsCount int

	articles := make([]Article, 0)
	links := make([]string, 0)

	//# Visit index pages and collect article links.
	c := colly.NewCollector()

	c.OnHTML(".b-022__block--description", func(e *colly.HTMLElement) {
		link := fmt.Sprintf("https://plus.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))
		links = append(links, link)
	})

	c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846?page=%d", pageNumber))

	//# Visit individual articles and collect article data.
	c = colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	//## Episode teaser
	c.OnHTML(".field-perex", func(e *colly.HTMLElement) {
		teaser = e.ChildText("p")
		// Try div element, somtimes there is empty <p> and <div> with content.
		// https://radiozurnal.rozhlas.cz/odpoledni-publicistika-emisni-povolenky-dovoz-ukrajinskeho-obili-kostra-8975193
		if teaser == "" {
			teaser = e.ChildText("div")
		}
		if teaser == "" {
			log.Println("Empty teaser", episode)
		}
	})

	//## Episode title
	c.OnHTML(".content", func(e *colly.HTMLElement) {
		episode = e.ChildText("h1")

		if strings.Contains(episode, "Polední") {
			time = "12:10"
			episode = strings.Replace(episode, "Polední publicistika: ", "", 1)
		} else {
			time = "18:10"
			episode = strings.Replace(episode, "Odpolední publicistika: ", "", 1)
		}

		date = convertDate(e.ChildText(".node-block__block--date"))

		topicsCount = countTopics(episode)
	})

	//## Episode moderator (may be missing, see `.factbox` parsing bellow when we use alternative strategy)
	c.OnHTML(".node-block--authors", func(e *colly.HTMLElement) {

		replacer := strings.NewReplacer("autor: ", "", "autoři:", "")

		moderatorText := strings.Split(e.Text, ",")[0]
		moderatorText = strings.TrimSpace(replacer.Replace(moderatorText))

		splitedModeratorText := strings.Split(moderatorText, " ")

		modedaratorFistName, moderatorLastName := splitedModeratorText[0], splitedModeratorText[1]

		moderator = Person{
			FirstName: modedaratorFistName,
			LastName:  moderatorLastName,
		}
	})

	//## Episode guests (+ moderator backup strategy when node-block--athors is missing)
	c.OnHTML(".factbox", func(e *colly.HTMLElement) {

		//### Scrape moderator
		// Use this strategy only when `moderator` is missing.
		// See https://radiozurnal.rozhlas.cz/poledni-publicistika-piratsky-sjezd-pavel-v-dnipru-knizni-festival-v-lipsku-8982816
		moderator2 = strings.Split(e.Text, ":")[0]
		replacer2 := strings.NewReplacer("Hosté", "", "Hosty", "", "byli", "", "je", "", "byl", "", "jsou", "")
		moderator2 = replacer2.Replace(moderator2)
		moderator2 = normalizeNames(moderator2)

		//### Scrape guests
		guests = make([]Guest, 0)
		guestsTexts := e.ChildTexts("li")

		// No <li> present so we try <p>.
		// See https://radiozurnal.rozhlas.cz/poledni-publicistika-piratsky-sjezd-pavel-v-dnipru-knizni-festival-v-lipsku-8982816
		if strings.TrimSpace(strings.Join(guestsTexts, "")) == "" {
			guestsTexts = e.ChildTexts("p")
			// The first paragraph may be "Hosty byly ..."
			if strings.Contains(guestsTexts[0], "Host") {
				guestsTexts = guestsTexts[1:]
			}
		}
		if strings.TrimSpace(strings.Join(guestsTexts, "")) == "" {
			logger.Println("Empty guests", episode)
		}

		replacer := strings.NewReplacer(".", " ", ";", " ")

		for i, g := range guestsTexts {
			guestsTexts[i] = replacer.Replace(g)
		}

		for _, g := range guestsTexts {
			if len(g) > 0 { // May be empty <li>!
				fields := strings.Fields(g)
				guests = append(guests, Guest{Person: Person{FirstName: fields[0], LastName: strings.ReplaceAll(fields[1], ",", "")}, Function: strings.Join(fields[2:], " ")})
			}
		}
	})

	//## Get all episodes
	for _, link := range links {
		c.Visit(link)
		// This will not work if async is allowed!
		if moderator.LastName == "" || moderator.FirstName == "" {
			split := strings.Split(moderator2, " ")
			fname, lname := split[0], split[1]
			moderator.LastName = lname
			moderator.FirstName = fname
		}
		article := NewArticle(show, episode, date, time, link, teaser, moderator, guests, topicsCount)
		articles = append(articles, article)
	}

	return articles
}
