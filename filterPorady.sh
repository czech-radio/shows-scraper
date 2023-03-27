#!/bin/bash

DATE=`date +%Y-%m-%d`
FILENAME=$DATE_porady_schedule.tsv
echo "date	since	title	moderators	description" > ${FILENAME}

cat *.csv | grep False | grep "Interview Plus"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}
cat *.csv | grep False | grep "Pro a proti"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}
cat *.csv | grep False | grep "Dvacet minut"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}
#cat *.csv | grep False | grep "rozhovory"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}


mv ${FILENAME} /root/irozhlas-scraper-geneea-output/publicistika/
rm *.csv
