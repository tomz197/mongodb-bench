#!/bin/bash

echo "Building MongoDB Benchmark Tool..."
go build -o mongodb-bench ./cmd/mongodb-bench

if [ $? -eq 0 ]; then
    echo "Build successful! Run with ./mongodb-bench -queries example-queries.json"
else
    echo "Build failed."
fi 