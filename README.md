# libkflow - library for sending kflow

libkflow is a library for generating and sending kflow records to Kentik
written in Go and also providing a simple C API.

## Go Usage

Create a libkflow `Config` and pass it to one of the `NewSenderWith*`
methods, along with the appropriate device identifier and an error
channel which may be polled for asynchronous errors.

    errors := make(chan error, 100)
    config := libkflow.NewConfig("email", "token", "program", "1.0.0")
    s, err := libkflow.NewSenderWithDeviceID(0, errors, config)
	s.Send(&flow.Flow{
		Ipv4SrcAddr: src,
		Ipv4DstAddr: dst,
	})

## C Usage

[demo.c](c/demo.c) provides an example of correct API usage and will
send one flow record to the test server which is expected to be
running and listening on 127.0.0.1:8999.

## Mock Server

A mock server that accepts API, flow, and metrics requests is available
for testing purposes:

    go install github.com/kentik/libkflow/cmd/server && bin/server

A variety of command line parameters are available, by default it will
listen on 127.0.0.1:8999 and accept some default values for
authentication and device identification. A full demonstration program
written in Go is available:

[main.go](cmd/demo/main.go)

## Build

Building libkflow requires GNU make and the Go toolchain. Set GOPATH to
the directory containing this file and invoke make:

    $ export PATH=$PWD
    $ make

The default make target will build `libkflow.a` and create a distribution
`libkflow-$VERSION-$TARGET.tar.gz` file containing the compiled library,
test server binary, demo.c, and kflow.h.

All built artifacts will be placed in `$CURDIR/out/$TARGET` where $CURDIR
is the current working directory and $TARGET is `$GOOS-$GOARCH`, for
example `linux-amd64`.

## Building library for internal use cases

Currently only supported when run from macOS.

Internal use cases require a version of libkflow targetting a variety of
architectures laid out in a specific directory structure. 

First, checkout the desired version from git and ensure it is tagged correctly. 
```shell
git tag -a $VERSION
```

First, ensure that the musl cross compilers are installed. This can be done
with the command:

```shell
brew install filosottile/musl-cross/musl-cross
```

Then, run the following command to build the libraries:

```shell
make internal-libs
```

This will build the libraries for all architectures and place them in the
`./bin/libs` directory. Continue steps from the internal use cases documentations.