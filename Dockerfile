FROM golang:1.15-buster as build-env

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

FROM scratch
COPY --from=build-env /app/server /server
COPY --from=build-env /app/web /web

ENV APP_ENV=production

EXPOSE 8080
ENTRYPOINT ["/server"]
