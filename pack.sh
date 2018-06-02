#!/usr/bin/bash

cd blade

rm -rf ../bin
mkdir -p ../bin/darwin
mkdir -p ../bin/linux
mkdir -p ../bin/windows

CGO_ENABLED=0 go build -ldflags '-w -s' -o ../bin/darwin/blade && upx ../bin/darwin/blade
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o ../bin/linux/blade && upx ../bin/linux/blade
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o ../bin/windows/blade && upx ../bin/windows/blade
