@echo off
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

set y=%date:~0,4%
set m=%date:~5,2%
set d=%date:~8,2%
set h=%time:~0,2%
set mi=%time:~3,2%
set s=%time:~6,2%

go build -o ../release/log%m%%d%%h%%mi% ../example/main.go
