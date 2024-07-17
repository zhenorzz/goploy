#!/bin/bash

OS=`uname`

echo "Change version number? [Version number/N]";
read x

if [ "$OS" = 'Darwin' ]; then
  CMD=gsed
else
  CMD=sed
fi

if [[ $x =~ ^[1-9].[0-9]+.[0-9]+$ ]]
then
  $CMD -i -e "s/const appVersion = \"[0-9].[0-9]\+.[0-9]\+\"/const appVersion = \"$x\"/g" cmd/server/main.go
  $CMD -i -e "s/@version [0-9].[0-9]\+.[0-9]\+/@version $x/g" cmd/server/main.go
  $CMD -i -e "s/'[0-9].[0-9]\+.[0-9]\+'/'$x'/g" database/goploy.sql
  $CMD -i -e "s/GOPLOY_VER=v[0-9].[0-9]\+.[0-9]\+/GOPLOY_VER=v$x/g" docker/Dockerfile
  $CMD -i -e "s/\"version\": \"[0-9].[0-9]\+.[0-9]\+\"/\"version\": \"$x\"/g" web/package.json
fi

echo "Build web? [Y/N]";
read x

if [ "$x" == Y ] || [ "$x" == y ]
then
    cd web
    npm run build
    cd ..
fi


echo "env GOOS=linux GOARCH=amd64  go build -o goploy cmd/server/main.go";
env GOOS=linux GOARCH=amd64 go build -o goploy cmd/server/main.go
echo "goploy done";

echo "env GOOS=linux GOARCH=arm64 go build -o goploy_arm64 cmd/server/main.go";
env GOOS=linux GOARCH=arm64 go build -o goploy_arm64 cmd/server/main.go
echo "goploy_arm64 done";

echo "env GOOS=darwin GOARCH=arm64 go build -o goploy_arm64.mac cmd/server/main.go";
env GOOS=darwin GOARCH=arm64 go build -o goploy_arm64.mac cmd/server/main.go
echo "goploy_arm64.mac done";

echo "env GOOS=darwin go build -o goploy.mac cmd/server/main.go";
env GOOS=darwin go build -o goploy.mac cmd/server/main.go
echo "goploy.mac done";

echo "env GOOS=windows GOARCH=amd64 go build -o goploy.exe cmd/server/main.go";
env GOOS=windows GOARCH=amd64 go build -o goploy.exe cmd/server/main.go
echo "goploy.exe done";

