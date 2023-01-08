# Host 服务模块

## IMPL
为Host服务模块的具体实现 围绕interface.go中Service接口进行编程

```
http
 |
Host Service (interface impl)
 |
impl(存储基于MySQL实现)

```

Host Service定义 并把实现编写完成, 使用方式有多种用途:
+ 用于内部模块调用, 基于他封装更高一层的业务逻辑, 比如发布服务
+ Host Service对外暴露: http协议(暴露给用户)
+ Host Service对外暴露: Grpc(暴露给内部服务)
+ ...