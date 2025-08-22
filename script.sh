#!/bin/bash

URLS=(
  "http://localhost:3002/posts"
  "http://localhost:3002/comments"
  "http://localhost:3002/todos"
)

COUNT=100

for i in $(seq 1 $COUNT)
do
  (
    URL=${URLS[$((i % ${#URLS[@]}))]}
    curl -s "$URL" > /dev/null
  ) &
done

wait
echo "Finished $COUNT parallel requests (round robin)."