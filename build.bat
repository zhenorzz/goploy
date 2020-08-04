statik -f -src=.\web\dist\
::go build linux
SET GOOS=linux
SET GOARCH=amd64
go build -o goploy main.go

::go build mac
SET GOOS=darwin
SET GOARCH=amd64
go build -o goploy.mac main.go

::go build windows
SET GOOS=windows
SET GOARCH=amd64
go build -o goploy.exe main.go
