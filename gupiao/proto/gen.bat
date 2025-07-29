protoc --go_out=./go --go-grpc_out=./go .\msg.proto
protoc --cpp_out=./cpp msg.proto
