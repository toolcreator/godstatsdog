FROM golang:alpine AS build
WORKDIR /go/src/godstatsdog
COPY . .
RUN go build

FROM alpine:latest
COPY --from=build go/src/godstatsdog/godstatsdog /usr/local/bin
ENTRYPOINT ["godstatsdog"]
