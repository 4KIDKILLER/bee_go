# bee_go

## 状态码
`200` 成功
`500-599` 服务器异常
`600-699` 业务流程错误
`700-799` 数据库错误

### go mod 包下载超时解决办法
```shell
  go env -w GOPROXY=https://goproxy.cn,direct
```
 这行命令会永久修改你的 Go 环境变量。direct 表示如果代理找不到，会直接去 GitHub 拉取，作为备选。 