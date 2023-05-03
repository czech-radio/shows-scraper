package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Person struct {
	FirstName string
	LastName  string
}

type Guest struct {
	Person
	Function string
}

type Article struct {
	Show      string
	Episode   string
	Date      string
	Time      string
	Link      string
	Teaser    string
	Moderator Person
	Guests    []Guest
}

type ShowName string

func NewArticle(show string, episode string, date string, time string, link string, teaser string, moderator Person, guests []Guest) Article {
	article := Article{}
	article.Show = show
	article.Episode = episode
	article.Date = date
	article.Time = time
	article.Link = link
	article.Teaser = teaser
	article.Moderator = moderator
	article.Guests = guests

	return article
}

func sortByDate(articles []Article) {
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

func (article *Article) PrettyPrint() {
	// return fmt.Sprintf("Pořad: %s\nEpizoda: %s\nDatum: %s\nObsah: %s\nLink : %s\n\n", article.Show, article.Episode, article.Date, article.Teaser, article.Link)

	fmt.Println(article.Show)
	fmt.Println(article.Date)
	fmt.Println(article.Time)
	fmt.Println(article.Link)
	fmt.Println(article.Teaser)
	fmt.Println(article.Moderator.FirstName, article.Moderator.LastName)
	for _, guest := range article.Guests {
		fmt.Println("*", guest.FirstName, guest.LastName, guest.Function)
	}

	fmt.Println("-------")
}

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

// GetRozhovoryEpisodes gets *Hlavní zprávy - rozhovory, komentáře* episodes.
func GetRozhovoryEpisodes(pageNumber int) []Article {

	show := "Hlavní zprávy - rozhovory a komentáře"
	var teaser, episode, date, time string
	var moderator Person
	var guests []Guest

	articles := make([]Article, 0)
	links := make([]string, 0)

	// Visit index pages and collect article links.
	// ------------------------------------------------------------------------- //
	c := colly.NewCollector()

	c.OnHTML(".b-022__block--description", func(e *colly.HTMLElement) {
		link := fmt.Sprintf("https://plus.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))
		links = append(links, link)
	})

	c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846?page=%d", pageNumber))

	// Visit individual articles and collect article data.
	// ------------------------------------------------------------------------- //
	c = colly.NewCollector()

	c.OnHTML(".field.field-perex", func(e *colly.HTMLElement) {
		teaser = e.ChildText("p")
	})

	c.OnHTML(".content", func(e *colly.HTMLElement) {
		episode = e.ChildText("h1")

		if strings.Contains(episode, "Polední") {
			time = "12:10"
		} else {
			time = "18:10"
		}

		date = convertDate(e.ChildText(".node-block__block--date"))
	})

	c.OnHTML(".node-block--authors", func(e *colly.HTMLElement) {
		moderatorText := e.ChildTexts("a")
		splitedModeratorText := strings.Split(moderatorText[0], " ")
		modedaratorFistName, moderatorLastName := splitedModeratorText[0], splitedModeratorText[1]
		moderator = Person{
			FirstName: modedaratorFistName,
			LastName:  moderatorLastName,
		}
	})

	c.OnHTML(".factbox", func(e *colly.HTMLElement) {
		guests = make([]Guest, 0)
		guestsText := strings.TrimSpace(e.ChildText("li"))
		for _, person := range strings.Split(guestsText, ";") {
			fields := strings.Fields(person)
			guests = append(guests, Guest{Person: Person{FirstName: fields[0], LastName: strings.ReplaceAll(fields[1], ",", "")}, Function: strings.Join(fields[2:], " ")})
		}
	})

	for _, link := range links {
		c.Visit(link)
		article := NewArticle(show, episode, date, time, link, teaser, moderator, guests)
		articles = append(articles, article)
	}

	return articles
}

// ------------------------------------------------------------------------- //

var (
	versionFlag     bool
	buildTime       string
	sha1GitRevision string
	versionGitTag   string
)

func main() {

	numPages := flag.Int("p", 1, "Number of pages to download.")
	flag.BoolVar(&versionFlag, "v", false, "Print application version and exit.")
	flag.Parse()

	// Must be at top!
	if versionFlag {
		fmt.Printf("Version: %s %s %s\n", versionGitTag, buildTime, sha1GitRevision)
		os.Exit(0)
	}

	articles := make([]Article, 0)

	for i := 0; i < *numPages; i++ {
		articles = append(articles, GetRozhovoryEpisodes(i)...)
	}

	// Sort articles in-place.
	sortByDate(articles)

	for _, article := range articles {
		article.PrettyPrint()
	}
}
