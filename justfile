# gomux task runner — run `just` to list recipes.

set shell := ["bash", "-euo", "pipefail", "-c"]

# List available recipes
default:
    @just --list --unsorted

# Build the gomux binary (./gomux, already gitignored)
build:
    go build -o gomux .

# Run gomux with optional arguments, e.g. `just run create --help`
run *args:
    go run . {{ args }}

# Run all tests, including real tmux integration tests
test *args:
    go test ./... {{ args }}

# Fast tests only: skips anything that starts tmux
test-unit:
    go test -short ./...

# Only the tmux integration tests, verbose and uncached
test-integration *args:
    go test ./internal/render/ -run TestTmuxIntegration -count=1 -v {{ args }}

# Tests for one roadmap milestone fixture, e.g. `just test-fixture single-window`
test-fixture name:
    go test ./internal/render/ -count=1 -v -run '/{{ name }}$'

# Rewrite golden files from the current generator output (review the diff before committing!)
update-golden:
    go test ./internal/render/ -run TestBashGolden -update
    git status --short -- internal/render/testdata

# Format all Go files
fmt:
    gofmt -w .

# Everything that should pass before a commit: formatting, vet, staticcheck, tests
check:
    #!/usr/bin/env bash
    set -euo pipefail
    unformatted="$(gofmt -l .)"
    if [[ -n "$unformatted" ]]; then
        echo "gofmt needed on:"; echo "$unformatted"; exit 1
    fi
    go vet ./...
    if command -v staticcheck >/dev/null; then
        staticcheck ./...
    else
        echo "note: staticcheck not installed, skipping (go install honnef.co/go/tools/cmd/staticcheck@latest)"
    fi
    go test ./...
