protoc -I=. --go_out=./apps --go_opt=module="github.com/Murphychih/cmdb/apps" --go-grpc_out=./apps --go-grpc_opt=module="github.com/Murphychih/cmdb/apps" ./apps/resource/pb/resource.proto

--proto_path: 指定了要去哪个目录中搜索import中导入的和要编译为.go的proto文件，可以定义多个
在根目录下执行操作

-I: 指定文件当前搜索目录