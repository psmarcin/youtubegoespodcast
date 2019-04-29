dev:
	GO111MODULE=on realize start --name="ygp-api"

dependencies:
	go get .

test: dependencies
	go test ./...

build: dependencies
	go build main.go

deploy:
	now

deploy_ci:
	echo $NOWSHTOKEN
	echo ${NOWSHTOKEN}
	now -t ${NOWSHTOKEN} --no-verify

alias:
	now alias

alias_ci:
	now alias -t ${NOWSHTOKEN} --no-verify

deploy_prod: deploy alias

deploy_prod_ci: deploy_ci alias_ci
