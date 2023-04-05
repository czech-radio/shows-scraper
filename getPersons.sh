#!/bin/bash
#

DATE="$1"

. ~/.env

activate(){
  # export GENEEA_API_KEY=xxxx
  cd geneea
  python3.10 -m venv .venv
  . .venv/bin/activate
  pip install -e .[dev]

}

getGeneea() {
 cro.geneea --input "$1" --type analysis --format json
}

deactivate(){
  deactivate
}

activate || exit 1

while read i; do
  echo "$i" | awk -F':' '{print $2}' > /tmp/input.txt
  getGeneea "/tmp/input.txt" || exit 1
done < /tmp/geneea_inputs.txt

deactivate || exit 1


exit 0
