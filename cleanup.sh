#!/bin/bash
DATE=`date +%Y-%m-%d`
FILENAME=${DATE}_publicistika.tsv

cp ${FILENAME} /root/irozhlas-scraper-geneea-output/publicistika/
rm *.csv
rm *.tsv
