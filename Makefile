GOPATH = $(shell echo $$GOPATH)

run:
	@ echo $(GOPATH)

install:
	@ CGO_ENABLED=0 go build -a -o $(GOPATH)/bin/kgen