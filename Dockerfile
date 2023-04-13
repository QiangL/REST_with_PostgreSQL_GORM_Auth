# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ae86-auth

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /ae86-auth /ae86-auth

EXPOSE 8080

ENTRYPOINT ["/ae86-auth"]
