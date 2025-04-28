module github.com/trysourcetool/sourcetool/examples/go/simple-create

go 1.24.0

toolchain go1.24.2

require github.com/trysourcetool/sourcetool-go v0.0.0-00010101000000-000000000000

require (
	github.com/gofrs/uuid/v5 v5.3.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/trysourcetool/sourcetool-go => ../../../sdk/go
