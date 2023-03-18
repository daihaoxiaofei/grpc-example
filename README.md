# gRPC的例子

###

    go语言调用 google.golang.org/grpc 包实现grpc的例子    

### 前置条件

    protoc 编译器下载
        https://github.com/protocolbuffers/protobuf/releases/latest
        选择类似 protoc-22.2-win64.zip 的文件 解压出 protoc.exe 放到环境变量中

    gRPC:
        go get -u github.com/golang/protobuf/proto
    
    插件 go专用的protoc生成器
        go get -u github.com/golang/protobuf/protoc-gen-go
    
    使用 protoc-gen-go 内置的gRPC插件生成gRPC代码:
        protoc --go_out=plugins=grpc:. protos/hello.proto

### 本例中知识点:

    基本传输
    双向传输
    服务发现与负载均衡
    拦截器
    附加内容传输
    鉴权
    接收到系统关闭信号时等待所有服务完成后关闭 除非停电暴力关机 不会临时终止服务
    
    目录 etcd-grpc 下是使用etcd做负载均衡的示列
    
