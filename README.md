# publicistika-scraper

*Program stáhne data o pořadu Rozhovory komentáře z webu Českého rozhlasu*

[![main](https://github.com/czech-radio/publicistika-scraper/actions/workflows/main.yml/badge.svg)](https://github.com/czech-radio/publicistika-scraper/actions/workflows/main.yml) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/518b8ee5b79240e78d3b955beb19d393)](https://app.codacy.com/gh/czech-radio/publicistika-scraper/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

## Build

```bash
git clone git@github.com:czech-radio/publicistika-scraper.git
cd publicistika-scraper
go build
```

## Usage

```bash
./publicistika-scraper -p [počet_stran]
```

- `p` Počet stran ke stažení, výchozí 1.
