.PHONY: test build dev curl

NAME := "fc"

test:
	go test -cover -v `glide nv`

build:
	go build -o $(NAME)
