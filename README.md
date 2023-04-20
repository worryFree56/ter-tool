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




|-- src
    |-- .DS_Store
    |-- .gitignore
    |-- .wuliu.yml
    |-- README.md
    |-- go.mod
    |-- go.sum
    |-- main.go
    |-- ter-tool
    |-- wuliu
    |-- build
    |-- cmd
    |   |-- doc.go
    |   |-- root.go
    |   |-- account
    |   |   |-- account.go
    |   |-- gpt
    |   |   |-- gpt.go
    |   |-- server
    |       |-- server.go
    |-- config
    |   |-- config.go
    |   |-- server.go
    |-- filelog
    |   |-- 202304200947.log
    |-- openapi
    |   |-- chat.go
    |   |-- chat_stream.go
    |   |-- chat_stream_reader.go
    |   |-- error.go
    |-- tool
        |-- aes.go
        |-- file.go
        |-- tool.go
