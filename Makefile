build: *.go
	GOARCH=amd64 GOOS=linux go build -o build/cLC -ldflags "-s -w"
	GOARCH=amd64 GOOS=windows go build -o build/cLC.exe -ldflags "-s -w"
