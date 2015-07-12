echo off

:: List of available target OSs and architectures:
:: https://golang.org/doc/install/source#environment

go generate

set distRoot=dist
mkDir %distRoot%
cd %distRoot%

::::::::::::::::::::::::::::::::

set goos=windows
set goarch=amd64

mkdir %goos%-%goarch%
cd %goos%-%goarch%

echo Building %goos%-%goarch%
go build ../..

::::::::::::::::::::::::::::::::

set goos=windows
set goarch=386

cd ..
mkdir %goos%-%goarch%
cd %goos%-%goarch%

echo Building %goos%-%goarch%
go build ../..

::::::::::::::::::::::::::::::::

set goos=linux
set goarch=386

cd ..
mkdir %goos%-%goarch%
cd %goos%-%goarch%

echo Building %goos%-%goarch%
go build ../..

::::::::::::::::::::::::::::::::

set goos=linux
set goarch=amd64

cd ..
mkdir %goos%-%goarch%
cd %goos%-%goarch%

echo Building %goos%-%goarch%
go build ../..