#!/bin/bash

set -e

go generate
go build

(while true; do 
	FREQ="$((80 + RANDOM % 400))"
	echo a $FREQ 3
	sleep 2
	FREQ="$((200 + RANDOM % 300))"
	echo b $FREQ 0.5
	sleep 0.5
	FREQ="$((180 + RANDOM % 400))"
	Y="$((FREQ/10))"
	echo a $FREQ 1 $Y
	sleep $(echo "scale=3;1 + $((RANDOM % 4000))/1000.0" | bc -l)
done) | ./soundofdata
