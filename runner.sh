#!/bin/bash
set -e

go test

# build the elm site
cd elm-site-source
elm make src/Main.elm --output ../main.js
cd ..
mv main.js public/

# move the binary so it doesn't clutter this directory
echo Starting the server
ls -al
go build *.go && mv main /tmp/blahblahmain && /tmp/blahblahmain
