module github.com/alpstable/gidari/examples/websocket

go 1.19

replace github.com/alpstable/gidari => ../../

require (
	github.com/alpstable/gidari v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.17.0
	google.golang.org/protobuf v1.30.0
)

require golang.org/x/time v0.3.0 // indirect
