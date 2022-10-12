#!/bin/bash
st=true
restart="1"
while [ $st ]
do
go run cmd/web/server.go = status
echo $status
    if [[ 1 = $status ]]; then
        echo "isItTrue"
        continue
    else
        break 1
    fi
done