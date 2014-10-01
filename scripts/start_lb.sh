#!/bin/bash

PORT=`cat config.json | grep Port | awk -F'\"' '{ print $4 }'`

if [ $PORT = "443" ]; then
	echo "nohup ./scalarm_load_balancer > log 2> log &" | sudo sh
else
	nohup ./scalarm_load_balancer > log 2> log &
fi
