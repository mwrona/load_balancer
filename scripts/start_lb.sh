#!/bin/bash

PORT=`cat config.json | grep Port | awk -F'\"' '{ print $4 }'`

if [ -z "$PORT" ]; then
	PORT=443
fi

if [ $PORT = "443" ]; then
	echo "nohup ./scalarm_load_balancer &" | sudo sh
else
	nohup ./scalarm_load_balancer
fi
