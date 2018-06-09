PACKAGE=solver

GOPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GOBIN=$(GOPATH)/bin
COVERAGEOUTPUT=coverage.out
COVERAGEHTML=coverage.html

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

run: build
	$(GOPATH)/bin/solver
