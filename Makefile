PACKAGE=solver

GOPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GOBIN=$(GOPATH)/bin
COVERAGEOUTPUT=coverage.out
COVERAGEHTML=coverage.html
IMAGENAME="kwkoo/$(PACKAGE)"
VERSION="0.1"

.PHONY: build clean test coverage run debug image container
build:
	@echo "Building..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(PACKAGE) $(PACKAGE)/cmd/$(PACKAGE)

clean:
	rm -f $(GOPATH)/bin/$(PACKAGE) $(GOPATH)/pkg/*/$(PACKAGE).a $(GOPATH)/$(COVERAGEOUTPUT) $(GOPATH)/$(COVERAGEHTML)

test:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test $(PACKAGE)

coverage:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test $(PACKAGE) -cover -coverprofile=$(GOPATH)/$(COVERAGEOUTPUT)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover -html=$(GOPATH)/$(COVERAGEOUTPUT) -o $(GOPATH)/$(COVERAGEHTML)
	open $(GOPATH)/$(COVERAGEHTML)

run:
	@GOPATH=$(GOPATH) go run $(GOPATH)/src/$(PACKAGE)/cmd/$(PACKAGE)/main.go

debug:
	@GOPATH=$(GOPATH) go run $(GOPATH)/src/$(PACKAGE)/cmd/$(PACKAGE)/main.go -debug true

image:
	docker build --rm -t $(IMAGENAME):$(VERSION) $(GOPATH)

runcontainer:
	docker run --rm -it --name $(PACKAGE) -p 8080:8080 $(IMAGENAME):$(VERSION)