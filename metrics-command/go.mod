module metrics-command

go 1.21.2

require example/metrics_events v1.0.0

require github.com/felixge/httpsnoop v1.0.3 // indirect

require (
	github.com/EventStore/EventStore-Client-Go v1.0.2
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
	google.golang.org/grpc v1.35.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace example/metrics_events => ../metrics_events
