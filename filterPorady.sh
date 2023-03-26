#!/bin/bash

DATE=`date +%Y-%m-%d`
FILENAME=porady_schedule_$DATE.tsv
echo "id	station	date	since	till	duration	repetition	title	moderators	description" > ${FILENAME}

cat *.csv | grep "Interview Plus" >> ${FILENAME}
cat *.csv | grep "Pro a proti" >> ${FILENAME}
cat *.csv | grep "Dvacet minut" >> ${FILENAME}
cat *.csv | grep "rozhovory" >> ${FILENAME}


