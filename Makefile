GO ?= go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
	COVERAGEDIR=$(CIRCLE_ARTIFACTS)/coverage
endif

ALL_PACKAGES = \
	bindings \
	endpoints \
	assured \

LDFLAGS = -ldflags '-X main.gitSHA=$(shell git rev-parse HEAD)'

MOCKERY = go run ./vendor/github.com/vektra/mockery/cmd/mockery/mockery.go

all: build test cover
install-deps:
	glide install
gen-mocks:
	for dir in $(MOCK_PACKAGES); do \
		echo "Generating Mocks under $$dir"; \
		$(MOCKERY) -dir=./$$dir/ -all -note "*DO NOT EDIT* Generated via mockery"; \
	done
	rm -rf ./mocks/error.go
build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build $(LDFLAGS) -v -o bin/go-rest-assured
fmt:
	find . -not -path "./vendor/*" -name '*.go' -type f | sed 's#\(.*\)/.*#\1#' | sort -u | xargs -n1 -I {} bash -c "cd {} && goimports -w *.go && gofmt -w -l -s *.go"
test:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	for dir in $(ALL_PACKAGES); do \
		$(GO) test -v ./$$dir -race -cover -coverprofile=$(COVERAGEDIR)/$$dir.coverprofile; \
	done
cover:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	for dir in $(ALL_PACKAGES); do \
		coverfile=$$(echo $$dir | awk '{gsub(/\//, "."); print}') && \
		$(GO) tool cover -html=$(COVERAGEDIR)/$${coverfile}.coverprofile -o $(COVERAGEDIR)/$${coverfile}.html; \
	done
assert-no-diff:
	test -z "$(shell git status --porcelain)"
clean:
	$(GO) clean
	rm -f bin/go-rest-assured
	rm -rf coverage/
	rm -rf vendor/
