module github.com/planx-lab/planx-plugin-processors

go 1.25.3

replace github.com/planx-lab/planx-sdk-go => ../planx-sdk-go

replace github.com/planx-lab/planx-proto => ../planx-proto

require github.com/planx-lab/planx-sdk-go v0.0.0-20260703110224-1e062d65dc66

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/planx-lab/planx-proto v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
