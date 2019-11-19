#!/bin/bash
# based on https://stackoverflow.com/questions/28320134/how-can-i-list-all-tags-for-a-docker-image-on-a-remote-registry
# This script fetches the tags of the provided Docker Hub repository and stores the latest tag (except 'latest') to disk.

if [ $# -lt 1 ]
then
cat << HELP

latest-tag  --  Get newest tag (except 'latest') from the remote registry

EXAMPLE USAGE: 
./latest-tag.sh scalair/robokops-base

HELP

exit 1
fi

image="$1"

echo "Retrieving tags for ${image}..."
raw_tags=`wget -q https://registry.hub.docker.com/v1/repositories/${image}/tags -O -`

echo "Parsing and sorting tags..."
latest_tag=`echo ${raw_tags} | sed -e 's/[][]//g' -e 's/"//g' -e 's/ //g' | tr '}' '\n' | awk -F: '{print $3}' | sort -r | grep -v latest | head -1`

echo "Latest tag: ${latest_tag}"

echo -n "${latest_tag}" > ./latest_tag

echo "Saved in 'latest_tag' file"
