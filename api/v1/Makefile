# Makefile to generate client and server stub code for proto-files.

# Dependencies (tools must be available on the path):
# - protoc (Protocol Buffer v3)
# - protoc-gen-go (protoc compiler extension for Go)
#
# That is install Go 1.6 or higher, gRPC, Protocol Buffers v3 and protoc plugin 
# for Go according to: https://grpc.io/docs/quickstart/go/
# See also: https://github.com/protocolbuffers/protobuf/releases

PROTO_FILES=quake.proto quake_api.proto
PROTOS=../../..
INPUT=$(PROTOS)/quake/api/v1
GO_OUT=$(PROTOS)

# default target builds Go stub 
all: go_stub

# Generate Go stubs
# (more about Go plugin: https://github.com/golang/protobuf)
go_stub: quake.pb.go
quake.pb.go: $(PROTO_FILES)
	protoc -I$(PROTOS) $(INPUT)/*.proto --go_out=plugins=grpc:$(GO_OUT)	
