#!/bin/bash
#

DATE=$1
PORAD="$2"



getSchedule() {
  cd schedule
  python -m venv .venv
  . .venv/bin/activate
  pip install -e .
  cro.schedule --period D --date $DATE --stations plus,radiozurnal --output ..
  deactivate
}


grepPorad() {
  cd ..
  # filtruj nazev poradu a premieru
  cat Schedule_D${DATE}.csv | grep "${PORAD}"
}


getSchedule || exit 1
grepPorad || exit 1

exit 0
