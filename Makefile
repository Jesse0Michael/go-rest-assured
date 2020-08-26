COVERAGEDIR = .coverage
GITSHA = $(shell git rev-parse HEAD)
LDFLAGS = -ldflags '-X main.gitSHA=$(GITSHA)'

all: build test cover
dependencies: 
	go mod download
build:
	if [ ! -d bin ]; then mkdir bin; fi
	go build $(LDFLAGS) -v -o bin/go-rest-assured
docker-build:
	docker build -f ./build/Dockerfile -t go-rest-assured:$(GITSHA) .
fmt:
	go mod tidy
	gofmt -w -l -s *.go
assert-no-diff:
	test -z "$(shell git status --porcelain)"
test:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	go test -v ./pkg/... -cover -coverprofile=$(COVERAGEDIR)/assured.coverprofile
cover:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	go tool cover -html=$(COVERAGEDIR)/assured.coverprofile
clean:
	go clean
	rm -f bin/go-rest-assured
	rm -rf $(COVERAGEDIR)
	rm -rf vendor/
