FROM golang:1.19-alpine AS build

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
