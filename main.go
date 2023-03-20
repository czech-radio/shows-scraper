package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	//	"strings"

	"github.com/gocolly/colly/v2"
)

type Option func(c Clanek) Clanek

type Clanek struct {
	Title       string
	Date        string
	Description string
	Link        string
	// optional
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
	no := strings.Split(input, ".")
	return fmt.Sprintf("%02d. %s", no[0], no[1])

}

func sortByDate(clanky []Clanek) {
	sort.SliceStable(clanky, func(i, j int) bool {
		ci, cj := prependZero(clanky[i].Date), prependZero(clanky[j].Date)

		switch {
		case ci != cj:
			return ci > cj
		default:
			return ci > cj
		}
	})
}

func sortNumbers(clanky []Clanek) ([]Clanek, error) {
	var lastErr error
	sort.Slice(clanky, func(i, j int) bool {
		a, err := strconv.ParseInt(clanky[i].Date, 0, 4)
		if err != nil {
			lastErr = err
			return false
		}
		b, err := strconv.ParseInt(clanky[j].Date, 0, 4)
		if err != nil {
			lastErr = err
			return false
		}
		return a < b
	})
	return clanky, lastErr
}

func (clanek *Clanek) PrettyPrint() {
	fmt.Printf("Název: %s\nDatum: %s\nObsah: %s\nLink: %s\n\n\n", clanek.Title, clanek.Date, clanek.Description, clanek.Link)
}

//// optioanl fields ///////////////////////////////////////////////////

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

/////////////////////////////////////////////////////////////////////////

func main() {
	c := colly.NewCollector()
	clanky := make([]Clanek, 0)

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		nadpis := e.ChildText("h3")
		datum := e.ChildText(".b-022__timestamp")
		popis := e.ChildText("p")
		link := fmt.Sprintf("https://radiozurnal.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))

		if nadpis != "" {
			clanky = append(clanky, NewClanek(nadpis, datum, popis, link))

		}
	})

	/*
			c.OnRequest(func(r *colly.Request) {
		          //fmt.Println("Visiting", r.URL)
			})
	*/

	//fmt.Println("\nstahuji: Hlavni zprávy rozhovory a komentáře\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846")
	//fmt.Println("\nstahuji: Pro a proti\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/pro-a-proti-6482952")
	//fmt.Println("\nstahuji: Dvacet minut radiožurnálu\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/dvacet-minut-radiozurnalu-5997743")
	//fmt.Println("\nstahuji: Interview Plus\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/interview-plus-6504167")

	// sort clanky by date
	/*
	        sorted, err := sortNumbers(clanky)
		if err != nil {
			fmt.Errorf("Error sorting by date %s", err.Error())
		}
	*/

	sortByDate(clanky)

	for _, clanek := range clanky {
		clanek.PrettyPrint()
	}

}
