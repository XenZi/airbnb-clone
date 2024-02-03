module reservation-service

go 1.20

require example/saga v1.0.0
require (
	github.com/gocql/gocql v1.6.0
	github.com/pariz/gountries v0.1.6
	github.com/sony/gobreaker v0.5.0
	go.opentelemetry.io/otel v1.17.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/sdk v1.17.0
	go.opentelemetry.io/otel/trace v1.17.0
)

require (
	github.com/cilium/ebpf v0.11.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mattn/go-isatty v0.0.3 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
	golang.org/x/arch v0.0.0-20190927153633-4e8777c89be4 // indirect
	golang.org/x/exp v0.0.0-20230224173230-c95f2b4c22f2 // indirect
	golang.org/x/sys v0.12.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/go-delve/delve v1.21.2
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace example/saga => ../saga