FROM golang:1.12 as build-env

RUN go version

ENV GO111MODULE=on
ENV APP_ENV=production

WORKDIR /app

# Copy all files
ADD . /app
RUN make build

# Tests
RUN make test

FROM gcr.io/distroless/base
COPY --from=build-env /app /
ENV APP_ENV=production
EXPOSE 8080
CMD ["/main"]
