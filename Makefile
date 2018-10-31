build: *.go lambda/*.go
	@GOARCH=amd64 GOOS=linux go build -o build/cLC -ldflags "-s -w"
	@upx build/cLC > /dev/null
	@GOARCH=amd64 GOOS=windows go build -o build/cLC.exe -ldflags "-s -w"
	@upx build/cLC.exe > /dev/null

nightly: build
	@mv build/cLC build/cLC_nightly
	@mv build/cLC.exe build/cLC_nightly.exe