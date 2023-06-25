# shows-scraper

[![main](https://github.com/czech-radio/shows-scraper/actions/workflows/main.yml/badge.svg)](https://github.com/czech-radio/shows-scraper/actions/workflows/main.yml) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/518b8ee5b79240e78d3b955beb19d393)](https://app.codacy.com/gh/czech-radio/shows-scraper/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

**Program stáhne data o pořadu [Hlavní zprávy - rozhovory a komentáře](https://radiozurnal.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846) z webu Českého rozhlasu**.

**This application is used to crowle/scrape data for internal usage.** 
It uses [colly](http://go-colly.org/). The final dataset is made with `process.ipynb`.

See documentation [here](https://github.com/czech-radio/organization/blob/main/software/projects/scraping-shows.md).

## Build

```bash
git clone git@github.com:czech-radio/shows-scraper.git
cd shows-scraper
```

### Unix

```bash
...
```

### Windows

```powershell
.\build.ps1
```

## Usage

```powershell
.\shows-scraper.exe -p [page_count (default: 1)]
```

např.

```powershell
.\shows-scraper.exe -p 3
```

## Poznámky
