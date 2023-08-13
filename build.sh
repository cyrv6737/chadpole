#!/bin/bash
GOOS=linux go build -o chadpole -ldflags="-s -w" *.go