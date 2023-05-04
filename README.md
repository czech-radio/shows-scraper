# rozhovory-scraper

[![main](https://github.com/czech-radio/rozhovory-scraper/actions/workflows/main.yml/badge.svg)](https://github.com/czech-radio/rozhovory-scraper/actions/workflows/main.yml) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/518b8ee5b79240e78d3b955beb19d393)](https://app.codacy.com/gh/czech-radio/rozhovory-scraper/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

**Program stáhne data o pořadu [Hlavní zprávy - rozhovory a komentáře](https://radiozurnal.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846) z webu Českého rozhlasu**.

**This application is used to crowle/scrape data for internal usage.** 
It uses [colly](http://go-colly.org/) and is fast and accurate.

- [x] Název pořadu
- [x] Název epizody
- [x] Datum vysílání epizody (`YYYY-MM-DD`).
- [x] Čas vysílání epizody (polední = 12:10, odpolední = 18:10).
- [x] Webový odkaz epizody
- [ ] Počet témat epizody (získáme z názvu, napovídá kolik bylo hostů).
- [x] Popis/Teaser epizody
- [x] Moderátor epizody (jméno, příjmení)
- [x] Hosté epizody (jméno, příjmení, funkce/popis).

Python kód je v adresáři [python](./python)
Výstupy se ukládájí do adresáře [data](./data)

## Build

```bash
git clone git@github.com:czech-radio/rozhovory-scraper.git
cd rozhovory-scraper
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
.\rozhovory-scraper.exe -p [page_count (default: 1)]
```

např.

```powershell
.\rozhovory-scraper.exe -p 3
```
