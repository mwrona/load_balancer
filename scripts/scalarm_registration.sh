#!/bin/bash
#config
INFORMATION_SERVICE_URL_="localhost:11300"
INFORMATION_SERVICE_LOGIN_="scalarm"
INFORMATION_SERVICE_PASSWORD_="scalarm"
REMOTE_LOAD_BALANCER_ADDRESS="localhost:9000"
LOCAL_LOAD_BALANCER_ADDRESS="localhost:9000"
#script
if [ -z "$INFORMATION_SERVICE_URL" ]; then
    INFORMATION_SERVICE_URL=$INFORMATION_SERVICE_URL_
fi
if [ -z "$INFORMATION_SERVICE_LOGIN_" ]; then
    INFORMATION_SERVICE_LOGIN_=$INFORMATION_SERVICE_LOGIN_
fi
if [ -z "$INFORMATION_SERVICE_PASSWORD" ]; then
    INFORMATION_SERVICE_PASSWORD=$INFORMATION_SERVICE_PASSWORD_
fi  

curl -u $INFORMATION_SERVICE_LOGIN:$INFORMATION_SERVICE_PASSWORD --data "address=$REMOTE_LOAD_BALANCER_ADDRESS" http://$INFORMATION_SERVICE_URL/experiment_managers
curl -u $INFORMATION_SERVICE_LOGIN:$INFORMATION_SERVICE_PASSWORD --data "address=$REMOTE_LOAD_BALANCER_ADDRESS/storage" http://$INFORMATION_SERVICE_URL/storage_managers

curl -k --data "address=$INFORMATION_SERVICE_URL&name=InformationSerivce" https://$LOCAL_LOAD_BALANCER_ADDRESS/register
 