module github.com/trysourcetool/sourcetool/examples/go/customer-support

go 1.22

toolchain go1.23.5

replace github.com/trysourcetool/sourcetool-go => ../../../sdk/go

require github.com/trysourcetool/sourcetool-go v0.0.0-20240415000000-000000000000

require (
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/samber/lo v1.47.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
