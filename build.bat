::go build linux
SET GOOS=linux
SET GOARCH=amd64
go build -o goploy main.go

::go build windows
SET GOOS=windows
SET GOARCH=amd64
go build -o goploy.exe main.go
