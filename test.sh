#!/bin/bash

export TESTING_BROKER_PORT=55001

docker run -d --cidfile container_id --rm -it \
  -p $TESTING_BROKER_PORT:1883 \
  -p 9001:9001 \
  -v $PWD/mosquitto/testing_passwords:/passwords \
  -v $PWD/mosquitto/mosquitto.conf:/mosquitto.conf \
  -v $PWD/mosquitto/mosquitto_acl:/mosquitto_acl eclipse-mosquitto > /dev/null

while ! nc -z localhost $TESTING_BROKER_PORT; do
  sleep 0.1 # wait for 1/10 of the second before check again
done

go test

docker stop $(cat container_id) > /dev/null
rm container_id
