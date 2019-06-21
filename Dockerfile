# ------------------------------------------------------------------------------
# Build Phase
# ------------------------------------------------------------------------------
FROM golang:1.11 AS build

ADD . /go/src/github.com/sosedoff/cloudwatchlogs
WORKDIR /go/src/github.com/sosedoff/cloudwatchlogs

RUN \
  GOOS=linux \
  GOARCH=amd64 \
  CGO_ENABLED=0 \
  go build -o /cloudwatchlogs

# ------------------------------------------------------------------------------
# Package Phase
# ------------------------------------------------------------------------------

FROM alpine:3.6

RUN \
  apk update && \
  apk add --no-cache ca-certificates openssl wget && \
  update-ca-certificates

WORKDIR /app

COPY --from=build /cloudwatchlogs /app/cloudwatchlogs
COPY static /app/static
COPY templates /app/templates

EXPOSE 5555
CMD ["/app/cloudwatchlogs"]