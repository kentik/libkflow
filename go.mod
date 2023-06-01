module github.com/kentik/libkflow

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/jessevdk/go-flags v1.4.0
	github.com/kentik/kit/go/legacy v0.0.2
	github.com/robfig/cron v0.0.0-20160927164231-9585fd555638
	github.com/stretchr/testify v1.8.0
	github.com/tinylib/msgp v1.1.6
	zombiezen.com/go/capnproto2 v2.18.2+incompatible
)

go 1.13

replace zombiezen.com/go/capnproto2 => github.com/kentik/go-capnproto2 v0.0.0-20221123170947-5a835423a561
