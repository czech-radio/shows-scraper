#!/bin/bash

DATE=`date +%Y-%m-%d`
FILENAME=${DATE}_porady_schedule.tsv
echo "date	since	title	moderators	description" > ${FILENAME}

cat *.csv | grep False | grep "Interview Plus"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}
cat *.csv | grep False | grep "Pro a proti"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}
cat *.csv | grep False | grep "Dvacet minut"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}
cat *.csv | grep False | grep -i "speciál"  | awk -F'\t' '{print $3"\t"$4"\t"$8"\t"$9"\t"$10}' >> ${FILENAME}


cp ${FILENAME} /root/irozhlas-scraper-geneea-output/publicistika/
rm *.csv
