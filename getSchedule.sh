#!/bin/bash
#

DATE="$1"


activate(){
  cd schedule
  python3.10 -m venv .venv
  . .venv/bin/activate
  pip install -e .

}

getSchedule() {
 cro.schedule --period D --date "$1" --stations plus,radiozurnal --output ..
}

deactivate(){
  deactivate
}

activate || exit 1

for i in `cat /tmp/dates.txt | sort -n -r |  uniq`; do
  getSchedule "$i" || exit 1
done

deactivate || exit 1


exit 0
