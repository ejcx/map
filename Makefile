NAME:=map
VERSION := $(shell git describe --tags --always --dirty="-dev")

build: 
	go build -o $(NAME) .

print-%  : ; @echo $* = $($*)

.PHONY: build
.SILENT: build
