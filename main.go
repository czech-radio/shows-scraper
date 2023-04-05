package main

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/tidwall/gjson"
)

type Option func(c Article) Article

type Article struct {
	Title       string
	Date        string
	Description string
	Link        string
	// optional
	Teaser    string
	Show      string
	Time      string
	Moderator string
	Guests    string
}

/* not used
type Show struct {
	Station     string `json:"station"`
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Since       string `json:"since"`
	Till        string `json:"till"`
	Repetition  string `json:"repetition"`
}
*/

func NewArticle(title string, date string, description string, link string, options ...Option) Article {
	c := Article{}
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
	fmt.Printf("Pořad: %s\nNázev: %s\nDatum: %s\nObsah: %s\nLink : %s\n\n", article.Show, article.Title, article.Date, article.Description, article.Link)
}

////////// WIP call geneea

func deriveGuests(article Article) Article {

	url := "https://api.geneea.com/v3/analysis/?T=CRo-transcripts"
	apiKey := fmt.Sprintf("%s", os.Getenv("GENEEA_API_KEY"))

	body := []byte(fmt.Sprintf(`{"text":"%s"}`, article.Description))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header = http.Header{
		"Host":          {url},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("user_key %s", apiKey)},
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	/*
		var geneea_reply map[string]interface{}

		derr := json.NewDecoder(res.Body).Decode(&geneea_reply)
		if derr != nil {
			panic(derr)
		}

		fmt.Println(geneea_reply)
	*/

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	/*
		unescaped, err := UnescapeUnicodeCharactersInJSON(resBody)
		if err != nil {
			fmt.Printf("there was an error unescaping json: %s\n", err.Error())
		}*/
	data := gjson.Get(string(resBody), "entities")

	data.ForEach(func(key, value gjson.Result) bool {

		attrs := gjson.GetMany(value.String(), "stdForm", "type")

		if attrs != nil {
			println(attrs[1].String())

			if attrs[1].String() == "person" && attrs[0].String() != article.Moderator {
				article.Guests = fmt.Sprintf("%s;%s", attrs[0].String(), article.Guests)
				//article.Moderator = attrs[4].String()
				//fmt.Printf("%s, %s, %s, %s, %s\n",article.Title,article.Date,article.Time,article.Description,article.Guests)
			}
		}
		return true
	})

	return article

}

type Person struct {
	givenName   string
	familyName  string
	description string
}

////////// WIP call schedules

func UnescapeUnicodeCharactersInJSON(jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func getSchedules(article Article) Article {

	split := strings.Split(article.Date, "-")
	year, month, day := split[0], split[1], split[2]

	id := "plus.json"
	url := "https://api.rozhlas.cz/data/v2"
	url = fmt.Sprintf("%s/%s/%s/%s/%s/%s", url, "schedule/day", year, month, day, id)

	// TODO API GET call here

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	//fmt.Printf("client: got response!\n")
	//fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	unescaped, err := UnescapeUnicodeCharactersInJSON(resBody)
	if err != nil {
		fmt.Printf("there was an error unescaping json: %s\n", err.Error())
	}
	data := gjson.Get(string(unescaped), "data")
	data.ForEach(func(key, value gjson.Result) bool {

		attrs := gjson.GetMany(value.String(), "title", "since", "description", "repetition", "persons.#.name")
		if attrs[0].String() == article.Show {
			article.Time = strings.Split(attrs[1].String(), "T")[1]
			article.Description = attrs[2].String()
			article.Moderator = attrs[4].String()
			//fmt.Printf("%s, %s, %s, %s, %s\n",article.Title,article.Date,article.Time,article.Description,article.Guests)
		}
		return true
	})

	return article
}

//// optional fields ///////////////////////////////////////////////////

func AddShow(show string) Option {
	return func(c Article) Article {
		c.Show = show
		return c
	}
}

func (article *Article) AddTime(time string) {
	article.Time = time
}

func (article *Article) AddModerator(moderator string) {
	article.Moderator = moderator
}

func (article *Article) AddGuests(guests string) {
	article.Guests = guests
}

func (article *Article) AddTeaser(teaser string) {
	article.Teaser = teaser
}

/////////////////////////////////////////////////////////////////////////

var showName string

func convertDate(input string) string {
	s := strings.Split(input, " ")
	day, err := strconv.Atoi(strings.Split(s[0], ".")[0])
	if err != nil {
		log.Fatal("Couldn't get day from date")
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

/////////////////////////////////////////////////////////////////////////
//////        MAIN  /////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

func main() {

	noPages := flag.Int("p", 1, "Number of pages to download.")
	flag.Parse()

	c := colly.NewCollector()
	articles := make([]Article, 0)

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		nadpis := e.ChildText("h3")
		if nadpis != "" {
			datum := convertDate(e.ChildText(".b-022__timestamp"))
			popis := e.ChildText("p")
			link := fmt.Sprintf("https://radiozurnal.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))

			novyArticle := NewArticle(nadpis, datum, popis, link, AddShow(showName))
			articles = append(articles, novyArticle)

		}
	})

	/*
			c.OnRequest(func(r *colly.Request) {
		          //fmt.Println("Visiting", r.URL)
			})
	*/

	for i := 0; i < *noPages; i++ {
		showName = "Hlavní zprávy - rozhovory, komentáře"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846?page=%d", i))

		showName = "Pro a proti"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/pro-a-proti-6482952?page=%d", i))

		showName = "Dvacet minut Radiožurnálu"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/dvacet-minut-radiozurnalu-5997743?page=%d", i))

		showName = "Interview Plus"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/interview-plus-6504167?page=%d", i))
	}

	sortByDate(articles)

	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	// TODO: get persons from moderator fields
	for index, article := range articles {
		articles[index] = getSchedules(article)
	}

	// TODO: call Geneea to mod description here
	for index, article := range articles {
		articles[index] = deriveGuests(article)
	}

	// write the complete output
	writeCSV(fmt.Sprintf("%s_publicistika.tsv", today), articles)

}

func runScript(command string) {
	cmd := exec.Command(command)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))
}

func clearTmp(filename string) {
	err := os.Remove(filename) // remove a single file
	if err != nil {
		fmt.Println(err)
	}

}

func writeFile(filename string, text string) {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func writeCSV(filename string, articles []Article) {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)
	w.Comma = '\t'

	header := []string{"Pořad", "Datum", "Čas", "Moderátor", "Host", "Název", "Popis", "Odkaz"}
	w.Write(header)

	defer w.Flush()
	var data [][]string
	for _, article := range articles {
		row := []string{article.Show, article.Date, article.Time, article.Moderator, article.Guests, article.Title, article.Description, article.Link}
		data = append(data, row)
	}
	w.WriteAll(data)

}

func readCsvFields(filePath string, articles []Article) []Article {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	// weird [i] but it works
	for i, article := range articles {
		for _, row := range records {
			if article.Date == row[0] && article.Show == row[2] {
				articles[i].Time = row[1]
				articles[i].Moderator = row[3]
				articles[i].Guests = row[4]
			}
		}
	}

	return articles
}
