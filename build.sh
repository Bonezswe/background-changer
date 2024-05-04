#!usr/bin/bash

echo "Building..."

go build -o cbg.exe \
    -ldflags "-X main.version=1.0.0" \
    .

echo "Build complete!"