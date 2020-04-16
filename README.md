

# 项目说明
基于rancher/cli项目进行的改造，通过调用 rancher api操作 --> k8s --> docker
* 适用于 rancher2.x
* 主要依赖
    * rancher/types       // rancher2.x API库
    * rancher/norman      // rancher/types 的支撑
    * goframe             // server化框架

# 环境和依赖
## 环境
    golang ： go version go1.13.4 linux/amd64
    rancher : 2.3

## 第三方依赖

* 环境设置
所有依赖包安装，都先设置环境变量
```shell
export GO111MODULE="on"
export GOPROXY=https://goproxy.cn
```

* golang.org相关包
可在以下目录寻找，并clone到 ${GOPATH}/src/golang.org/ 对应目录下
```shell
https://github.com/golang
```

# 通信协议
Json
## 输入
参看每个接口具体的情况
```shell
{
    "key":"value",
    "key":"value",
    "key":"value",
}
```

## 输出
```shell
{
    code:int,   //=0:成功; 非0:失败（错误码还未详细定义）
    msg:string, //code=0时统一为ok; code!=0时为错误提示
    data:array, //输出数据，没有数据时为空数组
}
```

# 程序结构
## 配置
参看 config/config.yaml

## server
入口是 main.go
运行后会启动一个 http 服务器

## api
参看 api/


# 运行和调用
## 运行
```shell
go run main.go
```

## 调用
启动后，需要先调用login接口，成功后的后续调用，需传入login接口返回的token参数
* 接口文档暂未提供
* 返回的信息暂未规整，返回的东西比较多             


# TODO
* 业务逻辑
    - 特别是网络部分
* 错误码规整
* 日志

# 欢迎交流
