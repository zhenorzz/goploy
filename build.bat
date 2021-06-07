if %1 == with-web (
    cd web

    call npm run build

    cd ..
)

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
