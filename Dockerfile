FROM golang:1.16-buster as build-env

ENV GO111MODULE=on
ENV APP_ENV=production

WORKDIR /app

ADD go.mod go.sum Makefile /app/
RUN make dependencies


# Copy all files
ADD . /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN make build-raw

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

FROM scratch
COPY --from=build-env /app/server /server

ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV APP_ENV=production

EXPOSE 8080
ENTRYPOINT ["/server"]
