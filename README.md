# 背景说明
rancher对应api的cli项目，github.com/rancher/cli(2.x版本的Rancher API)
将改造成一个Server版的Rancher API项目

# 安装和运行

go get github.com/rancher/cli
# # 处理一系列依赖问题后

# # 编译
go build main.go

# # 问题1：
```
[root@iZm5efctez2mq4wk8wbhsyZ cli]# go build main.go
# github.com/rancher/cli/cmd
cmd/multiclusterapp.go:723:5: app.Wait undefined (type *"github.com/rancher/types/client/management/v3".MultiClusterApp has no field or method Wait)
cmd/multiclusterapp.go:724:5: app.Timeout undefined (type *"github.com/rancher/types/client/management/v3".MultiClusterApp has no field or method Timeout)

# 解决方法：
进入 ./cmd/multiclusterapp.go
将723、724两行屏蔽，再次编译
```


# 使用
```
./main login https://47.104.225.225 --skip-verify --token token-2x6mg:bnlgj5msvxbg9hk5xgr9t88t5l557knrw6rcvrlndg5xhvt9qv7cq5
# 其中，--token 后部分是从rancher登录过程获取的会话token
```

# # 执行示例：
```
[root@iZm5efctez2mq4wk8wbhsyZ cli]# ./main login https://47.104.225.225 --skip-verify --token token-hpv9d:bc8rfl4gnxmsl796lzlnhsxbwzn5mtzjw7hl6gq5kgbpzj6m6q42r4
NUMBER    CLUSTER NAME   PROJECT ID        PROJECT NAME   PROJECT DESCRIPTION
1         test           c-9ktxk:p-fnlq6   Default        Default project created for the cluster
2         test           c-9ktxk:p-p7g9n   System         System project created for the cluster
3         test1          c-sr5l6:p-ksk7z   System         System project created for the cluster
4         test1          c-sr5l6:p-vtr8l   Default        Default project created for the cluster
Select a Project:1
INFO[0002] Saving config to /root/.rancher/cli2.json
```

# 文档

https://rancher.com/docs/rancher/v2.x/en/cli/  官方使用帮助，挺有用
