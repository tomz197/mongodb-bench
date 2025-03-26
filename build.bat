@echo off
echo Building MongoDB Benchmark Tool...
go build -o mongodb-bench.exe .\cmd\mongodb-bench

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Run with .\mongodb-bench.exe -queries example-queries.json
) else (
    echo Build failed.
) 