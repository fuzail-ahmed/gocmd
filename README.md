# gocmd — Tiny CLI Todo (Go)

A minimal command-line todo app written in Go — simple JSON-backed store for learning CLI patterns.

## Features
- Add, list, read, modify, delete tasks
- Stores tasks in `tasks.json`
- Small and easy to extend (good for learning Go CLI, file I/O, and basic patterns)

## Requirements
- Go 1.20+ (or your project's Go version)
- (Optional) `gofmt`, `go vet`, `golangci-lint` for linting

## Install & Run
```bash
# clone
git clone https://github.com/fuzail-ahmed/gocmd.git
cd gocmd

# build
go build -o gocmd

# or run directly
go run .

# usage examples:
./gocmd add --title "Buy milk" --desc "2 liters"
./gocmd list
./gocmd read --id 1
./gocmd modify --id 1 --title "Buy milk and bread" --completed true
./gocmd delete --id 1
