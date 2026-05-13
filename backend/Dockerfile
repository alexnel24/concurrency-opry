FROM golang:1.24 AS build

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

COPY . /src
RUN go build -o OpryScrape ./cmd

FROM debian:bookworm

WORKDIR /app

RUN mkdir -p /app/data

COPY --from=build /src/OpryScrape /app/

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

CMD ["./OpryScrape"]