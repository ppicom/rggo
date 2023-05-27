#!/bin/bash

OSLIST="linux darwin windows"
ARCHLIST="amd64 arm arm64"

for os in ${OSLIST}; do
    for arch in ${ARCHLIST}; do
        if [[ "$os/$arch" =~ ^(windows/arm64|darwin/arm) ]]; then continue; fi

        echo Building binary for $os $arch
        mkdir -p releases/${os}/${arch}
        CG_ENABLED=0 GOOS=$os GOARCH=$arch go build -tags=inmemory \
            -o releases/${os}/${arch}/
    done 
done