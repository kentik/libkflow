module github.com/kentik/libkflow

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/jessevdk/go-flags v1.4.0
	github.com/kentik/kit/go/legacy v0.0.2
	github.com/robfig/cron v0.0.0-20160927164231-9585fd555638
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.8.0
	github.com/tinylib/msgp v1.1.6
	go.uber.org/goleak v1.1.10
	zombiezen.com/go/capnproto2 v2.18.2+incompatible
)

require (
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.17

replace zombiezen.com/go/capnproto2 => github.com/kentik/go-capnproto2 v0.0.0-20221123170947-5a835423a561
