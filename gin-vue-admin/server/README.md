## server项目结构

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o xray-admin *.go
tar.exe -czvf xray-admin.tar.gz xray-admin resource
