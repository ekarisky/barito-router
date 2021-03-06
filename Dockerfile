FROM golang:1.13-buster AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p bin && \
  go build -o bin/ ./...

FROM ubuntu:18.04

COPY --from=build /app/bin/barito-router /usr/bin/barito-router

RUN useradd -m -U -d /app app
USER app
ENTRYPOINT [ "/usr/bin/barito-router" ]
