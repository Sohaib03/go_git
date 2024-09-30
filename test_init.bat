@echo off
if exist "test" (
    rmdir /S /Q "test"
)

go run .\main.go init test

if exist "test" (
    tree /F /A "test"
)
