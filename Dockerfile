FROM golang:1.14 as build-env

ENV GO111MODULE=on
ENV APP_ENV=production

WORKDIR /app

# Copy all files
ADD . /app
RUN make build

FROM gcr.io/distroless/base
COPY --from=build-env /app /
ENV APP_ENV=production
EXPOSE 8080
CMD ["/server"]
