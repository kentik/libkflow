#!/usr/bin/env bash
set -e

VERSION=$(git describe --tags --always --dirty)
function build_lib {
  OS=$1
  ARCH=$2
  OUTPUT_PATH=$3
  CC=$4

  echo "Building libkflow for $OS/$ARCH, output: $OUTPUT_PATH"
  if [ -z "$CC" ]; then
    GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=1 go build -o "$OUTPUT_PATH" -buildmode=c-archive -ldflags "-s -w -X github.com/kentik/libkflow.Version=${VERSION}" ./cmd/libkflow
  else
    GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=1 CC="$CC" go build -o "$OUTPUT_PATH" -buildmode=c-archive -ldflags "-s -w -extldflags -static -X github.com/kentik/libkflow.Version=${VERSION}" ./cmd/libkflow
  fi
}

# Build for macOS
build_lib darwin arm64 bin/libs/aarch64/macos/libkflow.a
build_lib darwin amd64 bin/libs/x86_64/macos/libkflow.a

# Build for linux/musl
build_lib linux arm64 bin/libs/aarch64/linux/musl/libkflow.a aarch64-linux-musl-gcc
build_lib linux amd64 bin/libs/x86_64/linux/musl/libkflow.a x86_64-linux-musl-gcc

# Builds with issues, probably due to time since last update (measured in years)
# Preserved here in case someone wants to try to fix them
#build_lib linux arm bin/libs/arm/linux/musl/libkflow.a
#build_lib linux amd64 bin/libs/x86_64/linux/libkflow.a
#build_lib freebsd amd64 bin/libs/x86_64/freebsd/libkflow.a


