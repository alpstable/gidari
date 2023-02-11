module github.com/alpstable/gidari/examples/csvpb

go 1.19

replace github.com/alpstable/gidari => ../../


require (
	github.com/alpstable/csvpb v0.1.0
	github.com/alpstable/gidari v0.1.0
	golang.org/x/time v0.3.0
)

require google.golang.org/protobuf v1.28.1 // indirect
