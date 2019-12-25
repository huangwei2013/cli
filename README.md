

# 项目说明
基于rancher/cli项目进行的改造
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

## 实体层次和调用前置条件说明
- cluster：预先设置(创建好)，无须改动
- project：预先设置(创建好)，无须改动
    - 对下面的实体，调用时都要指定 project(通过入参 projectId)
- namespace
- pipeline
    - 对应CI的配置
- pipeline execution
    - 对应一次CI的执行，所以依赖对应的pipeline(通过入参 pipelineId)
- deployment
    - 对应CD
    - 生成的镜像image，会直接通过rancher server写入regitry
        - rancher 的 registry 服务，预先设置好
- workload
    - 一个很大的信息集合，主要包括当前 cluster 中的所有活动实体信息
        - pods（包含containers信息）都在其中
        - images的信息，rancher目前是通过 workload 这个大实体来传递给前端的
                    


# TODO
* 业务逻辑
    - 特别是网络部分
* 错误码规整
* 日志

# 欢迎交流