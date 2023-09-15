FROM golang:1.21-alpine AS build

LABEL org.opencontainers.image.source "https://github.com/jesse0michael/go-rest-assurede"
LABEL org.opencontainers.image.description "Docker image for the GO Rest Assured client to mock and validate your calls to REST API's"

# Copy project files
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
RUN mkdir /dir

# Fetch dependencies
RUN go mod download
COPY . .

# Fetch dependencies (go mod)
RUN go mod download
RUN go mod verify

# Build GO binary
ENV CGO_ENABLED=0
RUN go build -o go-rest-assured ./cmd/go-assured

# Build final application
FROM gcr.io/distroless/static AS app

COPY --from=build /go/src/go-rest-assured ./

ENTRYPOINT ["./go-rest-assured"]
