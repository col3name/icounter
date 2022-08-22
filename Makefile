.PHONY: download migrateup migratedown lint build buildWin run up down herokulogs

download:
	go mod download

build: download
	set GOARCH=amd64
	set GOOS=linux
	set CGO_ENABLED=0
	go build -o bin/unique cmd/main.go

buildWin: download
	set GOARCH=amd64
	set GOOS=windows
	set CGO_ENABLED=0
	go build -o bin/unique.exe cmd/main.go

buildMacos: download
	GOOS=darwin GOARCH=arm64 go build -o bin/unique cmd/main.go

lint:
	golangci-lint run