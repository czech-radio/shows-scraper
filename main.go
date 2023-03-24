package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	//	"strings"

	"encoding/csv"
	"log"
	"os"
	"os/exec"

	"github.com/gocolly/colly/v2"
)

type Option func(c Clanek) Clanek

type Clanek struct {
	Title       string
	Date        string
	Description string
	Link        string
	// optional
	Teaser    string
	Show      string
	Time      string
	Moderator string
	Guests    []string
}

func NewClanek(title string, date string, description string, link string, options ...Option) Clanek {
	c := Clanek{}
	c.Title = title
	c.Date = date
	c.Description = description
	c.Link = link

	for _, o := range options {
		c = o(c)
	}

	return c
}

func prependZero(input string) string {
	s := strings.Split(input, ".")
	no, err := strconv.Atoi(s[0])

	if err != nil {
		fmt.Println("Error converting string to int: " + err.Error())
	}

	return fmt.Sprintf("%02d. %s", no, s[1])

}

func sortByDate(clanky []Clanek) {
	sort.SliceStable(clanky, func(i, j int) bool {
		ci, cj := clanky[i].Date, clanky[j].Date

		switch {
		case ci != cj:
			return ci > cj
		default:
			return ci > cj
		}
	})
}

func (clanek *Clanek) PrettyPrint() {
	fmt.Printf("Pořad: %s\nNázev: %s\nDatum: %s\nObsah: %s\nLink: %s\n\n\n", clanek.Show, clanek.Title, clanek.Date, clanek.Description, clanek.Link)
}

//// optional fields ///////////////////////////////////////////////////

func AddShow(show string) Option {
	return func(c Clanek) Clanek {
		c.Show = show
		return c
	}
}

func AddTime(time string) Option {
	return func(c Clanek) Clanek {
		c.Time = time
		return c
	}
}

func AddModerator(moderator string) Option {
	return func(c Clanek) Clanek {
		c.Moderator = moderator
		return c
	}
}

func AddGuests(guests []string) Option {
	return func(c Clanek) Clanek {
		c.Guests = guests
		return c
	}
}

func AddTeaser(teaser string) Option {
	return func(c Clanek) Clanek {
		c.Teaser = teaser
		return c
	}
}

/////////////////////////////////////////////////////////////////////////

var showName string

func getSchedule(date string, porad string) {
	cmd := exec.Command("./getSchedule.sh", date, porad)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Run())
}

func convertDate(input string) string {
	s := strings.Split(input, " ")
	day := strings.Split(s[0], ".")[0]
	year := fmt.Sprintf("%d", time.Now().Year())
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
	return fmt.Sprintf("%s-%s-%s", year, mo, day)
}

func Split(r rune) bool {
	return r == '.' || r == ' '
}

/////////////////////////////////////////////////////////////////////////

func main() {

	noPages := flag.Int("p", 1, "Number of pages to download.")
	flag.Parse()

	c := colly.NewCollector()
	clanky := make([]Clanek, 0)

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		nadpis := e.ChildText("h3")
		if nadpis != "" {
			datum := convertDate(e.ChildText(".b-022__timestamp"))
			popis := e.ChildText("p")
			link := fmt.Sprintf("https://radiozurnal.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))

			novyClanek := NewClanek(nadpis, datum, popis, link, AddShow(showName))
			clanky = append(clanky, novyClanek)

		}
	})

	/*
			c.OnRequest(func(r *colly.Request) {
		          //fmt.Println("Visiting", r.URL)
			})
	*/

	for i := 0; i < *noPages; i++ {
		showName = "Hlavní zprávy, rozhovory a publicistika"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846?page=%d", i))

		showName = "Pro a proti"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/pro-a-proti-6482952?page=%d", i))

		showName = "Dvacet minut Radiožurnálu"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/dvacet-minut-radiozurnalu-5997743?page=%d", i))

		showName = "Interview Plus"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/interview-plus-6504167?page=%d", i))
	}

	sortByDate(clanky)

	for _, clanek := range clanky {
		getSchedule(clanek.Date, clanek.Show)
		clanek.PrettyPrint()
	}

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
