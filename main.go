package main


import (
  "fmt"
  "github.com/gocolly/colly/v2"
)


func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".b-022__block", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("href"))
                datum := e.ChildText(".b-022__timestamp")
                nadpis := e.ChildText("h3")
                link := fmt.Sprintf("https://radiozurnal.rozhlas.cz%s", e.ChildAttr("h3 a","href"))
                popis := e.ChildText("p")
                
                if(nadpis != ""){
                  fmt.Printf("Název: %s\nDatum: %s\nLink: %s\nPopis: %s\n\n\n",nadpis,datum,link,popis)
                }
	})

	c.OnRequest(func(r *colly.Request) {
          //fmt.Println("Visiting", r.URL)
	})

        fmt.Println("\nstahuji: Hlavni zprávy rozhovory a komentáře\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846")
        fmt.Println("\nstahuji: Dvacet minut radiožurnálu\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/dvacet-minut-radiozurnalu-5997743")
        fmt.Println("\nstahuji: Interview Plus\n-----------------------------------------------------------\n\n")
	c.Visit("https://plus.rozhlas.cz/interview-plus-6504167")
}
