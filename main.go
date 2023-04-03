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

	//"bytes"
	"github.com/gocolly/colly/v2"
	"net/http"
	//"github.com/mohae/struct2csv"
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
	Guests    string
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
		ci, cj := fmt.Sprintf("%s %s", clanky[i].Show, clanky[i].Date), fmt.Sprintf("%s %s", clanky[j].Show, clanky[j].Date)

		switch {
		case ci != cj:
			return ci > cj
		default:
			return ci > cj
		}
	})
}

func (clanek *Clanek) PrettyPrint() {
	fmt.Printf("Pořad: %s\nNázev: %s\nDatum: %s\nObsah: %s\nLink : %s\n\n", clanek.Show, clanek.Title, clanek.Date, clanek.Description, clanek.Link)
}

func callGeneea(input string) string {

	url := "https://api.geneea.com/v3/analysis/T:CRo-transcripts"

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//Handle Error
	}

	apiKey := fmt.Sprintf("%s", os.Getenv("GENEEA_API_KEY"))

	req.Header = http.Header{
		"Host":          {"https://api.geneea.com/v3/analysis/T:CRo-transcripts"},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprtinf("user_key %s", apiKey)},
	}

	res, err := client.Do(req)
	if err != nil {
		//Handle Error
	}

	return res

}

func getSchedules(date string, stationId string) {

	split := strings.Split(date, "-")
	year, month, day := split[0], split[1], split[2]
	id = ""
	url := "https://api.rozhlas.cz/data/v2"
	dayData := fmt.Sprintf("%s/%s/%s/%s/%s/%s", url, "schedule/day", year, month, day, stationId)
	fmt.Println(dayData)
}

//// optional fields ///////////////////////////////////////////////////

func AddShow(show string) Option {
	return func(c Clanek) Clanek {
		c.Show = show
		return c
	}
}

func (clanek *Clanek) AddTime(time string) {
	clanek.Time = time
}

func (clanek *Clanek) AddModerator(moderator string) {
	clanek.Moderator = moderator
}

func (clanek *Clanek) AddGuests(guests string) {
	clanek.Guests = guests
}

func (clanek *Clanek) AddTeaser(teaser string) {
	clanek.Teaser = teaser
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
		showName = "Hlavní zprávy - rozhovory, komentáře"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846?page=%d", i))

		showName = "Pro a proti"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/pro-a-proti-6482952?page=%d", i))

		showName = "Dvacet minut Radiožurnálu"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/dvacet-minut-radiozurnalu-5997743?page=%d", i))

		showName = "Interview Plus"
		c.Visit(fmt.Sprintf("https://plus.rozhlas.cz/interview-plus-6504167?page=%d", i))
	}

	sortByDate(clanky)

	clearTmp("/tmp/dates.txt")

	for _, clanek := range clanky {
		//clanek.PrettyPrint()
		//getSchedule(clanek.Date, clanek.Show)
		//writeDates("/tmp/dates.txt",fmt.Sprintf("%s; %s\n",clanek.Date, clanek.Show))
		writeFile("/tmp/dates.txt", fmt.Sprintf("%s\n", clanek.Date))
	}

	runScript("./getSchedule.sh")
	runScript("./filterPorady.sh")

	currentTime := time.Now()
	today := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))

	// get schedule fileds
	clanky = readCsvFields(fmt.Sprintf("%s_porady_schedule.tsv", today), clanky)
	//fmt.Printf("Clanek.Time=%s", enrichedClanky[0].Time)

	// TODO: get persons from moderator fields
	clearTmp("/tmp/geneea_inputs.txt")

	for index, clanek := range clanky {
		writeFile("/tmp/geneea_inputs.txt", fmt.Sprintf("%02d: %s\n", index, clanek.Guests))
	}
	runScript("./getPersons.sh")

	// write the complete output
	writeCSV(fmt.Sprintf("%s_publicistika.tsv", today), clanky)

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

func writeCSV(filename string, clanky []Clanek) {
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
	for _, clanek := range clanky {
		row := []string{clanek.Show, clanek.Date, clanek.Time, clanek.Moderator, clanek.Guests, clanek.Title, clanek.Description, clanek.Link}
		data = append(data, row)
	}
	w.WriteAll(data)

}

func readCsvFields(filePath string, clanky []Clanek) []Clanek {
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
	for i, clanek := range clanky {
		for _, row := range records {
			if clanek.Date == row[0] && clanek.Show == row[2] {
				clanky[i].Time = row[1]
				clanky[i].Moderator = row[3]
				clanky[i].Guests = row[4]
			}
		}
	}

	return clanky
}
