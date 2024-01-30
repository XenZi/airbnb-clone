module metrics_query

go 1.21.2

require (
	example/metrics_events v0.0.0-00010101000000-000000000000
	github.com/EventStore/EventStore-Client-Go v1.0.2
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/grpc v1.61.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace example/metrics_events => ../metrics_events
