#!/bin/bash
#

DATE=$1
PORAD="$2"


activate(){
  cd schedule
  python -m venv .venv
  . .venv/bin/activate
  pip install -e .

}

getSchedule() {
 cro.schedule --period D --date $1 --stations plus,radiozurnal --output ..
}

deactivate(){
  deactivate
}

grepPorad() {
  cd ..
  # filtruj nazev poradu a premieru
  cat Schedule_D${DATE}.csv | grep True | grep "${PORAD}"
}


activate || exit 1

for i in `cat /tmp/dates.txt | sort -n |  uniq`; do
  getSchedule $i || exit 1
done
deactivate || exit 1


#grepPorad || exit 1

exit 0
