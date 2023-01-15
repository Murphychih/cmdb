## grpc
protoc -I=. --go_out=./apps --go_opt=module="github.com/Murphychih/cmdb/apps" --go-grpc_out=./apps --go-grpc_opt=module="github.com/Murphychih/cmdb/apps" ./apps/host/pb/host.proto
