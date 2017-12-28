#!/bin/bash

# http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-extldflags "-static"' .
docker build -t jimmykarily/smart_pie .
