FROM golang:alpine AS build
WORKDIR /go/src/godstatsdog
COPY . .
RUN apk add --no-cache git && go get -d -v ./... && go build

FROM alpine:latest
COPY --from=build go/src/godstatsdog/godstatsdog /usr/local/bin
ENTRYPOINT ["godstatsdog"]
