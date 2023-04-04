package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/csv"
	//"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	//"strings"
	//"bytes"

	"github.com/gocolly/colly/v2"
	"net/http"
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

func processGuests(article Article) Article {

	url := "https://api.geneea.com/v3/analysis/T:CRo-transcripts"

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//Handle Error
	}

	apiKey := fmt.Sprintf("%s", os.Getenv("GENEEA_API_KEY"))

	req.Header = http.Header{
		"Host":          {url},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("user_key %s", apiKey)},
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		//Handle Error
	}

	b, err := io.ReadAll(res.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(string(b))

	/*

		        var persons []Person
			return json.Unmarshall([]byte(res), &persons)
	*/

	article.Guests = fmt.Sprintf("%s", string(b))
	return article
}

type Person struct {
	fullName string
	role     string
}

////////// WIP call schedules

func getSchedules(article Article) Article {

	split := strings.Split(article.Date, "-")
	year, month, day := split[0], split[1], split[2]
	//id "0" = radiozurnal
	//id "3" = plus

	id := "3"
	url := "https://api.rozhlas.cz/data/v2"
	url = fmt.Sprintf("%s/%s/%s/%s/%s/%s", url, "schedule/day", year, month, day, id)

	// TODO API GET call here

	fmt.Println(url)
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

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
	fmt.Println(resBody)

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

/*
func getSchedule(date string, porad string) {
	cmd := exec.Command("./getSchedule.sh", date, porad)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Run())
}
*/

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

	/*
		clearTmp("/tmp/dates.txt")

		for _, article := range articles {
			//article.PrettyPrint()
			//getSchedule(article.Date, article.Show)
			//writeDates("/tmp/dates.txt",fmt.Sprintf("%s; %s\n",article.Date, article.Show))
			writeFile("/tmp/dates.txt", fmt.Sprintf("%s\n", article.Date))
		}

		runScript("./getSchedule.sh")
		runScript("./filterPorady.sh")
	*/
	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	// get schedule fileds
	// articles = readCsvFields(fmt.Sprintf("%s_porady_schedule.tsv", today), articles)
	//fmt.Printf("Article.Time=%s", enrichedClanky[0].Time)

	// TODO: get persons from moderator fields
	// call Geneea to mod description here
	for index, article := range articles {
		articles[index] = getSchedules(article)
	}

	// TODO: get persons from moderator fields
	// call Geneea to mod description here
	/*
	           for index, article := range articles {
	   		article[index] = processGuests(article)
	   	}
	*/
	/*
		        clearTmp("/tmp/geneea_inputs.txt")

			for index, article := range articles {
				writeFile("/tmp/geneea_inputs.txt", fmt.Sprintf("%02d: %s\n", index, article.Guests))
			}
			runScript("./getPersons.sh")
	*/

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
