module github.com/trysourcetool/sourcetool/examples/go/hello-world

go 1.22

replace (
	github.com/trysourcetool/sourcetool-go => ../../../sdk/go
	github.com/trysourcetool/sourcetool/proto => ../../../proto
)

require (
	github.com/trysourcetool/sourcetool-go v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.7.0
)

require (
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/samber/lo v1.47.0 // indirect
	github.com/trysourcetool/sourcetool/proto v0.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
