module reservation-service

go 1.21.6

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
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/nats-io/nats.go v1.32.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace example/saga => ../saga
