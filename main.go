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
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gocolly/colly/v2"
	"github.com/tidwall/gjson"
)

type Option func(c Article) Article

type Person struct {
	Prijmeni string
	Jmeno    string
	Funkce   string
}

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
	Guests    []Person
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

func deriveModerator(article Article) Article {

	url := "https://api.geneea.com/v3/analysis/?T=CRo-transcripts"
	apiKey := fmt.Sprintf("%s", os.Getenv("GENEEA_API_KEY"))

	body := []byte(fmt.Sprintf(`{"text":"%s"}`, article.Moderator))

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
			//println(attrs[1].String())

			article.Moderator = fmt.Sprintf("%s", attrs[0].String())

			//if attrs[1].String() == "person" && len(strings.Split(attrs[0].String(), " ")) <= 3 {
			//article.Guests = fmt.Sprintf("%s;%s", attrs[0].String(), article.Guests)
			//article.Guests = attrs[0].String()

			//article.Moderator = attrs[4].String()
			//fmt.Printf("%s, %s, %s, %s, %s\n",article.Title,article.Date,article.Time,article.Description,article.Guests)
			//}
		}
		return true
	})

	return article

}

func deriveGuests(article Article) Article {

	url := "https://api.geneea.com/v3/analysis/?T=CRo-transcripts"
	apiKey := fmt.Sprintf("%s", os.Getenv("GENEEA_API_KEY"))

	body := []byte(fmt.Sprintf(`{"text":"%s"}`, article.Teaser))

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

	persons := make([]Person, 0)

	data.ForEach(func(key, value gjson.Result) bool {

		attrs := gjson.GetMany(value.String(), "stdForm", "type")

		if attrs != nil {
			//println(attrs[1].String())

			if attrs[1].String() == "person" && len(strings.Split(attrs[0].String(), " ")) >= 2 {
				guests := strings.Split(attrs[0].String(), " ")
				name := guests[0]
				surname := guests[1]
				persons = append(persons, Person{Jmeno: name, Prijmeni: surname, Funkce: ""})

				article.Guests = persons

				//article.Guests = attrs[0].String()

				//article.Moderator = attrs[4].String()
				//fmt.Printf("%s, %s, %s, %s, %s\n",article.Title,article.Date,article.Time,article.Description,article.Guests)
			}
		}
		return true
	})

	return article

}

////////// WIP call schedules

func UnescapeUnicodeCharactersInJSON(jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// nevraci stejne vysledky jako python skipt
func getSchedules(article Article, stationId string) Article {

	split := strings.Split(article.Date, "-")
	year, month, day := split[0], split[1], split[2]

	id := fmt.Sprintf("%s.json", stationId)
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
		if attrs[0].String() == article.Show && attrs[2].String() == "false" {
			if attrs[1].String() != "" {
				article.Time = strings.Split(attrs[1].String(), "T")[1]
			}

			if attrs[2].String() != "" {
				article.Description = attrs[2].String()
			}

			if attrs[4].String() != "[]" {
				article.Moderator = fmt.Sprintf("%s;%s", attrs[4].String(), article.Moderator)
			}
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
	//article.Guests = guests
}

func (article *Article) AddTeaser(teaser string) {
	article.Teaser = teaser
}

/////////////////////////////////////////////////////////////////////////

var showName string

func convertDate(input string) string {
	s := strings.Split(input, " ")
	log.Println(input)
	day, err := strconv.Atoi(strings.Split(s[0], ".")[0])
	if err != nil {
		log.Println(fmt.Sprintf("Couldn't get day from date: %s", err.Error()))
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

var articles []Article

// Hlavní zprávy - rozhovory, komentáře
func A(articles []Article, i int) []Article {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		nadpis := e.ChildText("h3")
		if nadpis != "" {
			datum := convertDate(e.ChildText(".b-022__timestamp"))
			popis := e.ChildText("p")
			link := fmt.Sprintf("https://plus.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))

			novyArticle := NewArticle(nadpis, datum, popis, link, AddShow(showName))
			articles = append(articles, novyArticle)

		}
	})

	showName = "Hlavní zprávy - rozhovory, komentáře"
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

	// call Geneea to fix moderators
	for index, article := range articles {
		articles[index] = deriveModerator(article)
	}

	//fmt.Println(articles)

	return articles

}

// Pro a proti
func B(articles []Article, i int) []Article {

	c := colly.NewCollector()

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

	showName = "Pro a proti"
	c.Visit(fmt.Sprintf("https://radiozurnal.rozhlas.cz/pro-a-proti-6482952?page=%d", i))

	// visit links
	c = colly.NewCollector()

	var teaser string
	c.OnHTML(".field.field-perex", func(e *colly.HTMLElement) {
		teaser = fmt.Sprintf(e.ChildText("p"))
	})

	for i, article := range articles {
		c.Visit(article.Link)

		// Define the separators as a function
		separators := func(c rune) bool {
			return c == '?' || c == '.' || c == ';'
		}

		// Split the string by the separators
		sentences := strings.FieldsFunc(teaser, separators)

		articles[i].Teaser = strings.Join(sentences[1:len(sentences)], ".")

	}

	// call Geneea to fix moderators
	for index, article := range articles {
		articles[index] = deriveModerator(article)
	}

	// call Geneea to fix guests
	for index, article := range articles {
		articles[index] = deriveGuests(article)
	}

	clearTmp("/tmp/dates.txt")
	for _, article := range articles {
		writeFile("/tmp/dates.txt", fmt.Sprintf("%s\n", article.Date))
	}
	runScript("./getSchedule.sh")
	runScript("./filterPorady.sh")

	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	articles = readCsvFields(fmt.Sprintf("%s_porady_schedule.tsv", today), articles)

	return articles
}

// Dvacet minut Radiožurnálu
func C(articles []Article, i int) []Article {

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		nadpis := e.ChildText("h3")
		if nadpis != "" {
			popis := e.ChildText("p")

			datum := convertDate(e.ChildText(".b-022__timestamp"))
			link := fmt.Sprintf("https://radiozurnal.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))

			novyArticle := NewArticle(nadpis, datum, popis, link, AddShow(showName))
			articles = append(articles, novyArticle)
		}
	})

	showName = "Dvacet minut Radiožurnálu"
	c.Visit(fmt.Sprintf("https://radiozurnal.rozhlas.cz/dvacet-minut-radiozurnalu-5997743?page=%d", i))

	c = colly.NewCollector()

	var teaser string
	c.OnHTML(".field.field-perex", func(e *colly.HTMLElement) {
		teaser = fmt.Sprintf(e.ChildText("p"))
	})

	for i, article := range articles {
		c.Visit(article.Link)
		articles[i].Teaser = teaser
	}

	clearTmp("/tmp/dates.txt")
	for _, article := range articles {
		writeFile("/tmp/dates.txt", fmt.Sprintf("%s\n", article.Date))
	}
	runScript("./getSchedule.sh")
	runScript("./filterPorady.sh")

	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	articles = readCsvFields(fmt.Sprintf("%s_porady_schedule.tsv", today), articles)

	/*
		// call Geneea to fix moderators
		for index, article := range articles {
			articles[index] = deriveModerator(article)
		}
	*/

	// call Geneea to fix guests
	for index, article := range articles {
		articles[index] = deriveGuests(article)

		if len(articles[index].Guests) >= 1 {
			last := articles[index].Guests[len(articles[index].Guests)-1]
			fmt.Println(last)
			articles[index].Guests = []Person{{Jmeno: last.Jmeno, Prijmeni: last.Prijmeni, Funkce: last.Funkce}}
		}

	}

	return articles

}

// Interview Plus
func D(articles []Article, i int) []Article {

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		nadpis := e.ChildText("h3")
		if nadpis != "" {
			popis := e.ChildText("p")

			datum := convertDate(e.ChildText(".b-022__timestamp"))
			link := fmt.Sprintf("https://plus.rozhlas.cz%s", e.ChildAttr("h3 a", "href"))

			novyArticle := NewArticle(nadpis, datum, popis, link, AddShow(showName))
			articles = append(articles, novyArticle)
		}
	})

	showName = "Interview Plus"
	c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/interview-plus-6504167?page=%d", i))

	c = colly.NewCollector()

	var teaser string
	c.OnHTML(".field.field-perex", func(e *colly.HTMLElement) {
		teaser = fmt.Sprintf(e.ChildText("p"))
	})

	for i, article := range articles {
		c.Visit(article.Link)
		articles[i].Teaser = teaser
	}

	clearTmp("/tmp/dates.txt")
	for _, article := range articles {
		writeFile("/tmp/dates.txt", fmt.Sprintf("%s\n", article.Date))
	}
	runScript("./getSchedule.sh")
	runScript("./filterPorady.sh")
	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	articles = readCsvFields(fmt.Sprintf("%s_porady_schedule.tsv", today), articles)

	/*
		// call Geneea to fix moderators
		for index, article := range articles {
			articles[index] = deriveModerator(article)
		}
	*/

	// call Geneea to fix guests
	for index, article := range articles {
		articles[index] = deriveGuests(article)

		if len(articles[index].Guests) >= 1 {
			last := articles[index].Guests[len(articles[index].Guests)-1]
			fmt.Println(last)
			articles[index].Guests = []Person{{Jmeno: last.Jmeno, Prijmeni: last.Prijmeni, Funkce: last.Funkce}}
		}

	}

	return articles

}

/////////////////////////////////////////////////////////////////////////
//////        MAIN  /////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

func main() {

	noPages := flag.Int("p", 1, "Number of pages to download.")
	flag.Parse()

	articlesA := make([]Article, 0)
	articlesB := make([]Article, 0)
	articlesC := make([]Article, 0)
	articlesD := make([]Article, 0)

	for i := 0; i < *noPages; i++ {
		articlesA = A(articlesA, i)
		articlesB = B(articlesB, i)
		articlesC = C(articlesC, i)
		articlesD = D(articlesD, i)
	}

	articles := make([]Article, 0)
	articles = append(articles, articlesA...)
	articles = append(articles, articlesB...)
	articles = append(articles, articlesC...)
	articles = append(articles, articlesD...)

	sortByDate(articles)
	fmt.Println(articles)

	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	// TODO: get persons from moderator fields
	/*
		for index, article := range articles {
			articles[index] = getSchedules(article, "radiozurnal")
			articles[index] = getSchedules(article, "plus")
		}
	*/
	/*
		clearTmp("/tmp/dates.txt")
		for _, article := range articles {
			writeFile("/tmp/dates.txt", fmt.Sprintf("%s\n", article.Date))
		}
		runScript("./getSchedule.sh")
		runScript("./filterPorady.sh")

		articles = readCsvFields(fmt.Sprintf("%s_porady_schedule.tsv", today), articles)
	*/

	runScript("./cleanup.sh")
	// write the complete output
	writeXLS(fmt.Sprintf("%s_publicistika.xls", today), articles)
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

func writeXLS(filename string, articles []Article) {

	f := excelize.NewFile()

	// Create a new sheet in the XLS file.
	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	// Define the column headers for the sheet.
	columns := map[string]string{
		"A1": "datum",
		"B1": "čas",
		"C1": "pořad",
		"D1": "moderátor",
		"E1": "příjmení",
		"F1": "jméno",
		"G1": "strana",
		"H1": "popis_funkce",
		"I1": "funkce",
		"J1": "popis_téma",
	}

	// Write the column headers to the sheet.
	for col, header := range columns {
		f.SetCellValue(sheetName, col, header)
	}

	// Write the data from the slice of structs to the sheet.
	for i, article := range articles {
		row := i + 2 // Add 2 to skip the header row.
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), article.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), article.Time)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), article.Show)

		if article.Moderator != "" {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), article.Moderator)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "n/a (auto)")
		}

		if len(article.Guests) > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), article.Guests[0].Prijmeni)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), article.Guests[0].Jmeno)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "n/a (auto)")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "n/a (auto)")
		}

		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "n/a (auto)")
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "n/a (auto)")

		if len(article.Guests) > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), article.Guests[0].Funkce)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), "n/a (auto)")

		}
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), article.Teaser)
	}

	// Save the XLS file.
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
}

func writeXML(filename string, articles []Article) {

	xmlData, err := xml.MarshalIndent(articles, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling XML:", err)
		return
	}

	// Save the XML data to a file
	err = ioutil.WriteFile(filename, xmlData, 0644)
	if err != nil {
		fmt.Println("Error writing XML file:", err)
		return
	}

	fmt.Println("XML file saved successfully!")

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
		row := []string{article.Show, article.Date, article.Time, article.Moderator, fmt.Sprintf("%s, %s, %s;", article.Guests[0].Prijmeni, article.Guests[0].Jmeno, article.Guests[0].Funkce), article.Title, article.Description, article.Link}
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
				//articles[i].Guests = row[4]
			}
		}
	}

	return articles
}
