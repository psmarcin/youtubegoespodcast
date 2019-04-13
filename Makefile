dependencies:
	go get .

test: dependencies
	go test ./...

build: dependencies
	go build main.go

deploy:
	now

alias:
	now alias

deploy_prod: deploy alias
