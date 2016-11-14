# libkflow - library for sending kflow

libkflow is a library for generating and sending kflow records to Kentik
written in Go and callable from C with a simple API.

## Usage

[demo.c](demo.c) provides an example of correct API usage and will send
one flow record to the test server which is expected to be running and
listening on 127.0.0.1:8999.

## Build

Building libkflow requires GNU make and the Go toolchain. Set GOPATH to
the directory containing this file and invoke make:

```
$ export PATH=$PWD
$ make
```

The default make target will build `libkflow.a` and create a distribution
`libkflow-$VERSION-$TARGET.tar.gz` file containing the compiled library,
test server binary, demo.c, and kflow.h.

All built artifacts will be placed in `$CURDIR/out/$TARGET` where $CURDIR
is the current working directory and $TARGET is `$GOOS-$GOARCH`, for
example `linux-amd64`.
