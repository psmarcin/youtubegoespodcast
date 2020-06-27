FROM golang:1.14 as build-env

ENV GO111MODULE=on
ENV APP_ENV=production

WORKDIR /app

# Copy all files
ADD . /app
RUN GOOS=linux GOARCH=amd64 make build

FROM gcr.io/distroless/base
COPY --from=build-env /app/server /server
COPY --from=build-env /app/web /web

ENV APP_ENV=production

EXPOSE 8080
ENTRYPOINT ["/server"]
