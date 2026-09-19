#!/usr/bin/bash

go build -o main

docker build -t venin/go-final-project .

docker run -d -p 7540:7540 venin/go-final-project