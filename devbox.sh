#!/usr/bin/env sh

build() {
    if [ -z "$OUTPUT" ]; then
        OUTPUT="build/ren"
    fi

    go build -tags "$TAGS" -o "$OUTPUT" ./cmd/ren
}

eval $@
