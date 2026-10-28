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

# 项目开发记录

### 20260929

- `bee_user`添加二级密码`private_key`字段
- `bee_file`表`status`字段修改备注

### 20261012

- `bee_file`添加`is_private`字段

### 20261028

- 添加`online.yaml`环境配置，`-env`支持设置`online`

### 打包命令

```text
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/bee_go ./cmd/main.go
```

1. CGO_ENABLED=0
禁用 CGO。Go 默认可能调用系统 C 库，交叉编译到 Linux 时 Mac 上找不到 Linux 的 C 库会报错。设为 0 后生成纯静态链接的二进制，不依赖目标系统的任何动态库，丢到 Ubuntu 上就能直接运行。这是交叉编译 Linux 程序最关键的一项。

2. GOOS=linux
指定目标操作系统为 Linux。这是交叉编译的核心——在 Mac 上编译出 Linux 能执行的 ELF 文件，而不是 Mac 的 Mach-O 文件。

3. GOARCH=amd64
指定目标 CPU 架构为 64 位 x86。绝大多数 Ubuntu 服务器是这个架构。如果服务器是 ARM（如部分云实例），要改成 arm64。

这三个环境变量写在命令前面，只对这一次 go build 生效，不会污染当前终端会话，也不会影响你 Mac 上正常的本地编译。
4. go build
Go 的编译命令，把源码编译成可执行文件。

5. -ldflags="-s -w"
传给链接器的参数，用于减小二进制体积：

-s：去掉符号表（symbol table）
-w：去掉 DWARF 调试信息
两者合计通常能让文件小 20%–30%。代价是无法用 gdb 等调试器调试该二进制，也无法拿到完整的 panic 堆栈符号。生产部署一般可以接受；如果还需要事后调试，就去掉这段。

6. -o ./build/bee_go
指定输出文件路径为 ./build/bee_go。-o 后面是输出，不是输入。注意 ./build/ 目录必须已存在，否则会报错，可先 mkdir -p ./build。

7. ./cmd/main.go
编译目标，即要编译的源文件。这里指定的是单个 main.go 文件。如果 cmd 目录下有多个 .go 文件（比如还依赖同包的其他文件），写成 ./cmd/main.go 可能编译不完整，更稳妥的写法是用包路径 ./cmd（前提是 cmd 是一个完整的 main 包）。

查看指定服务日志
```shell
sudo journalctl -u bee_go.service -n 80 --no-pager
```