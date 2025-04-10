set dotenv-load

export EDITOR := 'nvim'

files := 'src/*'

alias f := fmt
alias r := run

default:
  just --list

all: test lint forbid fmt-check

run *args:
	#!/bin/bash
	go run `fd .go ./src -E *_test.go` {{ args }}

test:
	go test -v ./src

fmt:
	golines -m 80 -w {{files}}
	just retab

fmt-check:
	gofmt -l .
	@echo formatting check done

forbid:
	./bin/forbid

lint:
  golangci-lint run {{files}}

retab:
	./bin/retab

dev-deps:
	brew install golangci-lint golines
