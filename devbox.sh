#!/usr/bin/env sh

test() {
    export CGO_ENABLED=1
    go test -race ./...
}

build() {
    if [ -z "$OUTPUT" ]; then
        OUTPUT="build/ren"
    fi

    go build -o "$OUTPUT" ./cmd/ren
}

eval $@
