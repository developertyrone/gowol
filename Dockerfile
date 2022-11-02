ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /web
WORKDIR /web

COPY src/go.mod .
COPY src/go.sum .
RUN go mod download

COPY src .
RUN go build -o ./app ./main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /web
WORKDIR /web
COPY --from=builder /web/app .
COPY --from=builder /web/static ./static
COPY --from=builder /web/templates ./templates
COPY --from=builder /web/config.json .

EXPOSE 8080

ENTRYPOINT ["/web/app"]