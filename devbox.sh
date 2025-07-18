#!/usr/bin/env sh

export MODULES="$(tr '\n'  ' ' < default_modules.txt)"

test() {
    if [ -z "$TAGS" ]; then
        TAGS="$MODULES"
    fi

    export CGO_ENABLED=1
    go test -tags "$TAGS" -race ./...
}

build() {
    if [ -z "$OUTPUT" ]; then
        OUTPUT="build/ren"
    fi

    if [ -z "$TAGS" ]; then
        TAGS="$MODULES"
    fi

    go build -tags "$TAGS" -o "$OUTPUT" ./cmd/ren
}

eval $@
