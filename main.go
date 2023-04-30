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

type Option func(c Article) Article

type Person struct {
	Prijmeni string
	Jmeno    string
	Funkce   string
}

type Article struct {
	Show        string
	Episode     string
	Date        string
	Description string
	Link        string
	Teaser      string
	Time        string
	Moderator   string
	Guests      []Person
}

type ShowName string

func NewArticle(show string, title string, date string, description string, link string, options ...Option) Article {
	article := Article{}
	article.Show = show
	article.Episode = title
	article.Date = date
	article.Description = description
	article.Link = link

	for _, option := range options {
		article = option(article)
	}

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
	fmt.Printf("Pořad: %s\nNázev: %s\nDatum: %s\nObsah: %s\nLink : %s\n\n", article.Show, article.Episode, article.Date, article.Description, article.Link)
}

func convertDate(input string) string {
	s := strings.Split(input, " ")
	day, err := strconv.Atoi(strings.Split(s[0], ".")[0])
	if err != nil {
		// log.Println(fmt.Sprintf("Couldn't get day from date: %s", err.Error()))
		return input

	}
	year := fmt.Sprintf("%s", s[2])
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
func GetRozhovoryEpisodes(i int) []Article {

	articles := make([]Article, 0)

	c := colly.NewCollector()

	// Find and visit all links to episodes.
	show := "Hlavní zprávy - rozhovory a komentáře"
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		episode := e.ChildText("h3")
		if episode != "" {
			date := convertDate(e.ChildText(".b-022__timestamp"))
			desc := e.ChildText("p")
			link := fmt.Sprintf("https://plus.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))
			article := NewArticle(show, episode, date, desc, link)
			articles = append(articles, article)

		}
	})

	c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846?page=%d", i))

	c = colly.NewCollector()

	cnt := 0
	var moderator, guests, reply, teaser string
	var persons []Person

	c.OnHTML(".field.field-perex", func(e *colly.HTMLElement) {
		teaser = fmt.Sprintf(e.ChildText("p"))
	})

	c.OnHTML(".factbox", func(e *colly.HTMLElement) {
		reply = fmt.Sprintf(e.ChildText("p"))
		reply2 := fmt.Sprintf(e.ChildText("li") + " ")

		inline := fmt.Sprintf("%s %s", reply, reply2)

		split := strings.Split(inline, ":")
		moderator = strings.ReplaceAll(split[0], "Hosty", "")
		moderator = strings.ReplaceAll(moderator, "Hosté", "")
		moderator = strings.ReplaceAll(moderator, "jsou", "")
		moderator = strings.ReplaceAll(moderator, "byli", "")
		moderator = strings.ReplaceAll(moderator, `"`, "")

		moderator = strings.TrimSpace(moderator)

		guests = split[1]
		guests = strings.ReplaceAll(guests, `"`, "")
		guests = strings.TrimSpace(guests)
		entries := strings.Split(guests, ";")

		persons = make([]Person, 0)

		for _, person := range entries {
			fields := strings.Fields(person)
			persons = append(persons, Person{Jmeno: fields[0], Prijmeni: strings.ReplaceAll(fields[1], ",", ""), Funkce: strings.Join(fields[2:len(fields)], " ")})
		}

		cnt++
	})

	for i, article := range articles {
		c.Visit(article.Link)
		articles[i].Moderator = moderator
		articles[i].Guests = persons
		articles[i].Teaser = teaser
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

	// sortByDate(articles)

	for _, item := range articles {
		fmt.Println(item.Date)
		fmt.Println(item.Time)
		fmt.Println(item.Show)
		fmt.Println(item.Episode)
		fmt.Println(item.Description)
		fmt.Println(item.Moderator)
		fmt.Println("----")
	}
}
