# ter-tool

## mac 

### x86
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./

### arm
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build main.go




## window 
### x86
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o wuliu ./

### arm
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build main.go