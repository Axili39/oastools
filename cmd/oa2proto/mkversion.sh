#!/bin/sh
tag=$(git describe --tags --always --long --dirty)
echo "package main ; const ( version=\"$tag\" )" | gofmt > version.go
