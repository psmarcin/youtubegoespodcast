NAME=youtubegoespodcast
PROJECT_ID=youtubegoespodcast
IMAGE_TAG=gcr.io/$(PROJECT_ID)/$(NAME)

dev:
	modd

lint-ci:
	gofmt -w -s .

lint: lint-ci
	golangci-lint run --fix

dependencies:
	go get ./...

test:
	go test ./...

build: dependencies build-raw

build-raw:
	go build ./cmd/server

debug:
	dlv debug --headless --listen=:2345 --log --api-version 2

setProject:
	gcloud config set project $(PROJECT_ID)

docker-build:
	docker build -t $(NAME) .

docker-run: build
	docker run -p 8080:8080 $(NAME)

deploy: setProject
	gcloud builds submit --tag $(IMAGE_TAG)

publish: setProject
	gcloud beta run deploy --image $(IMAGE_TAG) --allow-unauthenticated --timeout=10 --concurrency=100 --memory=128Mi --region=us-central1 --update-env-vars=APP_ENV=production,API_URL=https://ygp.psmarcin.dev/ --platform managed $(NAME)
