.PHONY: download migrateup migratedown lint build buildWin run up down herokulogs

download:
	go mod download

build: download
	export GOARCH=amd64
	export GOOS=linux
	export CGO_ENABLED=0
	go build -o bin/unique cmd/main.go

buildWin: download
	set GOARCH=amd64
	set GOOS=windows
	set CGO_ENABLED=0
	go build -o bin/unique.exe cmd/main.go

buildMacos: download
	GOOS=darwin GOARCH=arm64 go build -o bin/unique cmd/main.go

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out

lint:
	golangci-lint run