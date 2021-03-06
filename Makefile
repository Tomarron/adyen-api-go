SHELL:=/bin/bash

test:
	go test -parallel 5 -v ./...

verification:
	go vet ./...
	golint -set_exit_status ./...
	CGO_ENABLED=0 errcheck ./...