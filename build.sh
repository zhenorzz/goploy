#!/bin/bash

echo "Change version number? [Version number/N]";
read x

if [[ $x =~ ^[1-9].[0-9]+.[0-9]+$ ]]
then
  sed -i -e "s/const appVersion = \"[0-9].[0-9]\+.[0-9]\+\"/const appVersion = \"$x\"/g" cmd/server/main.go
  sed -i -e "s/'[0-9].[0-9]\+.[0-9]\+'/'$x'/g" model/sql/goploy.sql
  sed -i -e "s/GOPLOY_VER=v[0-9].[0-9]\+.[0-9]\+/GOPLOY_VER=v$x/g" docker/Dockerfile
fi

echo "Build web? [Y/N]";
read x

if [ "$x" == Y ] || [ "$x" == y ]
then
    cd web
    npm run build
    cd ..
fi


echo "Building goploy";

env GOOS=linux go build -o goploy cmd/server/main.go

env GOOS=darwin go build -o goploy.mac cmd/server/main.go

env GOOS=windows go build -o goploy.exe cmd/server/main.go


