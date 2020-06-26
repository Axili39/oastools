#!/bin/sh
tag=$(git describe --always --long --dirty)
echo "package main ; const ( version=\"$tag\" )" | gofmt -w version.go
