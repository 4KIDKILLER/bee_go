# 文件上传接口性能问题报告与优化方案

## 1. 背景与结论

部署环境是一台老笔记本，4 核心/4 线程，4GB 内存。前端使用 10 条 HTTP 并发上传约 200 张图片，并且前端已经限制大于 10MB 的文件不会发起上传；上传到约 100 张时内存被打满并导致服务宕机。

从当前代码看，宕机的核心原因不是“200 张 10MB 图片被一次性保存到内存”，而是上传链路存在两处内存放大：

1. 上传接口使用 `r.ParseMultipartForm(32 << 20)`，而单文件上限是 10MB。因为 multipart 内存阈值高于单文件大小，10MB 以内的文件会优先进入内存，再由服务层从内存复制到磁盘。
2. 每次上传成功后，服务层立即启动一个不受限制的缩略图 goroutine。缩略图生成会把压缩图片完整解码成像素矩阵，实际内存占用通常远大于文件大小。当上传速度高于缩略图处理速度时，缩略图 goroutine 会积压，最终撑爆 4GB 内存。

前端 10MB 限制可以减少超大请求风险，但不能限制图片解码后的像素内存，也不能限制后端缩略图任务数量。在这台机器上，当前实现不适合承接“10 并发 + 近 200 张图片 + 自动缩略图”的批量上传场景。推荐优先做三个改动：降低 multipart 内存阈值或改为流式上传、限制缩略图并发、增加图片像素上限校验。

## 2. 当前上传调用链

当前文件上传入口：

- `pkg/controller/file_controller.go`
  - `POST /upload`
  - `http.MaxBytesReader`
  - `r.ParseMultipartForm`
  - `r.FormFile("file")`
  - 调用 `UploadFileService`

上传保存和缩略图生成：

- `pkg/service/file_service.go`
  - 创建上传目录
  - 生成文件 ID
  - `io.Copy(dst, file)` 保存原图
  - 写入 `bee_file` 表
  - `go func() { utils.ImageCompression(...) }()` 异步生成缩略图

图片处理：

- `pkg/utils/tools.go`
  - `imaging.Open` 完整解码图片
  - `imaging.Resize` 生成缩略图图像
  - `imaging.Save` 重新编码保存

简化流程：

```text
HTTP request
  -> MaxBytesReader(10MB)
  -> ParseMultipartForm(32MB)
       -> 10MB 以内文件保存在内存
  -> FormFile
  -> io.Copy 到磁盘
  -> INSERT bee_file
  -> 启动一个新的缩略图 goroutine
       -> imaging.Open 完整解码原图
       -> imaging.Resize 创建新图像
       -> imaging.Save 编码落盘
       -> UPDATE bee_file.file_thumb_path
```

## 3. 主要问题清单

| 优先级 | 问题 | 代码位置 | 影响 |
| --- | --- | --- | --- |
| P0 | 缩略图 goroutine 不限并发 | `pkg/service/file_service.go:117` | 上传越快，后台图片解码任务越多，内存峰值不可控，是本次 OOM 的最大嫌疑点 |
| P0 | 图片处理完整解码原图 | `pkg/utils/tools.go:50`、`pkg/utils/tools.go:87` | 文件小于 10MB 不代表内存小；一张 24MP 图片解码后约 96MB，多个任务并发会迅速耗尽内存 |
| P1 | multipart 解析阈值设置过高 | `pkg/controller/file_controller.go:62` | 10MB 以内文件被保存在内存中，10 并发时会产生至少约 100MB 的请求体内存峰值 |
| P1 | 没有上传接口级并发限制 | `pkg/controller/file_controller.go:50` | `net/http` 会为请求创建 goroutine，前端并发、多个用户或异常客户端都可以把后端并发推高 |
| P1 | 没有图片像素上限和格式白名单 | `pkg/controller/file_controller.go:75` 后仅检查文件大小 | 小体积高像素图片、恶意图片或非图片文件会进入保存和异步处理流程 |
| P2 | 请求体上限等于文件上限 | `pkg/controller/file_controller.go:59` | multipart 边界、文件字段和普通字段也占请求体大小，接近 10MB 的合法图片可能被误判超限 |
| P2 | HTTP Server 未设置超时 | `pkg/infrastructure/server/http_server.go:108` | 慢速上传或异常连接可能长期占用 goroutine、连接和内存 |
| P2 | DB 连接池对低内存机器偏高 | `config/config-prod.yaml:9`、`config/config-prod.yaml:10` | Go 服务和 MySQL 同机时，10 个打开连接和 10 个空闲连接会增加内存压力 |
| P3 | 缩略图更新错误日志判断有 bug | `pkg/service/file_service.go:122` 到 `pkg/service/file_service.go:124` | `thumbErr` 失败时可能不记录，排查缩略图失败会更困难 |

## 4. 内存峰值推导

### 4.1 multipart 阶段

当前代码：

```go
r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
r.ParseMultipartForm(32 << 20)
```

`ParseMultipartForm(maxMemory)` 的行为是：将 multipart 表单解析完整，文件内容在低于 `maxMemory` 时优先放在内存，更大的文件才会落到临时文件。

当前前端和后端都按 10MB 控制单文件大小，multipart 内存阈值是 32MB，所以几乎所有合法上传图片都会先进入内存。前端 10 并发时，仅请求体文件内容就可能形成：

```text
10 并发 * 10MB = 100MB+
```

这还不包含 multipart 元数据、表单字段、Go 对象开销、响应对象、JWT 解析、日志、数据库驱动缓冲等。

单看 100MB 通常不至于压垮 4GB 内存，但它会和缩略图任务叠加，形成更高峰值。

### 4.2 图片解码阶段

图片文件大小不是图片处理内存占用。前端限制大于 10MB 的文件不上传，只能约束压缩后的文件体积；JPEG/PNG 是压缩格式，处理前必须解码成像素矩阵。按 RGBA/NRGBA 估算：

```text
解码后内存约等于 宽 * 高 * 4 字节
```

常见照片的解码内存示例：

| 分辨率 | 像素数 | 单张解码后约占用 |
| --- | ---: | ---: |
| 4000x3000 | 12MP | 48MB |
| 6000x4000 | 24MP | 96MB |
| 8000x6000 | 48MP | 192MB |

当前 `ImageCompression` 中还会发生：

- `imaging.Open` 解码原图。
- `imaging.AutoOrientation(true)` 可能产生额外图像副本。
- `imaging.Resize` 创建缩略图目标图像。
- `imaging.Save` 编码输出，JPEG/PNG 编码也会产生额外临时内存和 CPU 压力。

因此一张小于 10MB 的图片，在生成缩略图时的瞬时内存可能是几十 MB 到数百 MB。

如果上传请求结束得很快，而缩略图生成较慢，代码会不断创建新的后台 goroutine。假设 24MP 图片每个缩略图任务瞬时占用 120MB 左右：

```text
10 个缩略图任务并发 约 1.2GB
20 个缩略图任务并发 约 2.4GB
30 个缩略图任务并发 约 3.6GB
```

这还没有计算 Go 运行时、MySQL、操作系统、文件缓存和 multipart 上传内存。对于 4GB 内存的机器，上传到约 100 张时宕机是符合当前实现特征的。

## 5. 为什么会在约 100 张时爆内存

前端始终保持 10 条上传并发。单个上传请求保存原图和写 DB 后就返回成功，但缩略图处理还在后台继续运行。

当缩略图生成速度低于上传完成速度时，系统状态会变成：

```text
前端继续补齐 10 个上传请求
  -> 后端持续接收和保存原图
  -> 每个成功上传都新增一个缩略图 goroutine
  -> 缩略图 goroutine 队列实际无限增长
  -> 多个 goroutine 同时解码大图
  -> 内存峰值持续上升
  -> GC 无法及时回收仍在使用的图片对象
  -> OOM 或系统开始疯狂换页，最终服务宕机
```

这也是为什么问题不是一开始就出现，而是上传一段时间后集中爆发。

## 6. 快速止血方案

这些方案改动小，适合先让老笔记本稳定跑起来。

### 6.1 前端上传并发从 10 降到 2 到 3

在 4 核心/4GB 内存机器上，10 路上传并发过高。建议先改成：

```text
上传 HTTP 并发：2 或 3
```

如果后端缩略图仍然不限并发，前端降并发只能缓解，不能根治。

### 6.2 将 multipart 内存阈值降低到 1MB

把：

```go
r.ParseMultipartForm(32 << 20)
```

改为：

```go
r.ParseMultipartForm(1 << 20)
```

这样大部分图片会进入临时文件，而不是完整保留在 Go 堆内。这个方案代码改动最小，但会增加一次临时文件读写，属于用磁盘换内存。

### 6.3 立即限制缩略图生成并发

最小可行方案是在服务层加一个全局信号量，让同一时刻最多只有 1 到 2 个缩略图任务执行。

对这台机器建议：

```text
缩略图 worker 数：1
上传接口并发：2 到 3
```

如果只做一个后端改动，优先做这个。

### 6.4 给 Go 进程设置内存目标

如果使用 Go 1.19 及以上版本，可以在启动服务时设置：

```shell
GOMEMLIMIT=1200MiB GOGC=50 ./bee_go -env prod
```

建议值需要考虑同机 MySQL、系统桌面环境和文件缓存。如果 MySQL 也运行在这台 4GB 机器上，Go 服务的目标内存建议先控制在 1.2GB 到 1.5GB。

注意：`GOMEMLIMIT` 只能让 GC 更积极，不能解决无限 goroutine 和图片解码并发失控。它是保护网，不是根治方案。

### 6.5 下调数据库连接池

`config/config-prod.yaml` 当前：

```yaml
maxOpenConns: 10
maxIdleConns: 10
```

低内存单机部署建议先改为：

```yaml
maxOpenConns: 4
maxIdleConns: 2
```

数据库不是本次 OOM 的主因，但同机 MySQL 会额外消耗内存，连接池不宜按服务器规格配置。

## 7. 推荐的后端优化方案

### 7.1 上传解析改为低内存或流式

推荐分两阶段做。

第一阶段，低风险改造：

- 保留 `ParseMultipartForm`。
- 将 `maxMemory` 从 32MB 降到 1MB 或更低。
- 将请求体上限改为 `maxFileSize + multipartHeadroom`。

示例：

```go
const maxFileSize int64 = 10 << 20
const multipartHeadroom int64 = 1 << 20

r.Body = http.MaxBytesReader(w, r.Body, maxFileSize+multipartHeadroom)
if err := r.ParseMultipartForm(1 << 20); err != nil {
    // 按错误类型区分 400 和 413
    return
}
```

第二阶段，推荐重构：

- 使用 `r.MultipartReader()` 逐 part 读取。
- 文件 part 直接 `io.CopyBuffer` 到最终文件。
- 用 `io.LimitedReader` 或自定义计数 reader 限制单文件大小。
- 普通字段只限制小尺寸并保存在内存。

这样上传过程不需要把完整文件放入 Go 堆，也不需要先写临时文件再复制到最终目录。

### 7.2 增加上传接口并发限制

后端必须有自己的限流，不应该完全相信前端并发设置。

建议配置：

```yaml
file:
  maxUploadConcurrency: 3
```

接口处理逻辑：

```go
select {
case uploadSem <- struct{}{}:
    defer func() { <-uploadSem }()
default:
    w.Header().Set("Retry-After", "3")
    fileController.writeError(w, http.StatusTooManyRequests, "上传任务繁忙，请稍后重试")
    return
}
```

在低内存机器上，返回 `429 Too Many Requests` 比让服务宕机更可控。

### 7.3 缩略图生成改为有界 worker 池

把当前的：

```go
go func() {
    compErr := utils.ImageCompression(dstPath, thumbDstPath, fileExt)
    ...
}()
```

改为：

```text
上传成功
  -> 投递 ThumbJob 到有界队列
  -> 固定 1 到 2 个 worker 消费队列
  -> worker 生成缩略图并更新 DB
```

建议配置：

```yaml
file:
  thumbWorkers: 1
  thumbQueueSize: 50
```

低配机器推荐 `thumbWorkers: 1`。这样最多只有 1 张图片在执行高内存解码，内存峰值会稳定很多。

队列满时有两种策略：

1. 阻塞当前上传，形成后端背压。
2. 返回 `202 Accepted` 或上传成功但标记缩略图待处理，由后台任务稍后补齐。

当前项目更简单的做法是阻塞或短暂等待队列，不建议无限创建 goroutine。

### 7.4 增加图片格式与像素上限校验

当前只校验 `fileHeader.Size > maxFileSize`，没有校验是否真的是图片，也没有校验像素数。

建议在保存后、缩略图前使用 `image.DecodeConfig` 读取图片尺寸：

```go
cfg, format, err := image.DecodeConfig(f)
if err != nil {
    return ErrInvalidImage
}
if cfg.Width*cfg.Height > maxImagePixels {
    return ErrImagePixelsTooLarge
}
```

建议参数：

```text
允许格式：jpg、jpeg、png
最大像素：20000000 到 24000000
```

对 4GB 内存机器，建议先用 20MP。超过上限的图片可以要求前端先压缩，或者后端单独走更慢、更严格的离线处理。

### 7.5 优化缩略图尺寸与算法

当前缩略图最大宽度或高度大致是 1280/720。对于相册列表预览，如果不是大图预览，可以降低到：

```text
列表缩略图：最长边 512 到 720
详情预览图：最长边 1280，可异步生成
```

还可以拆成两种图：

- thumb：小缩略图，快速生成，用于列表。
- preview：中等尺寸预览图，可后台慢慢生成。

这能显著降低编码后的文件体积和后续静态资源访问压力。

### 7.6 设置 HTTP Server 超时

当前 `http.Server` 只设置了 `Handler` 和 `Addr`。建议补充：

```go
return &http.Server{
    Handler:           middleware(mux),
    Addr:              config.Server.Port,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       5 * time.Minute,
    WriteTimeout:      1 * time.Minute,
    IdleTimeout:       60 * time.Second,
    MaxHeaderBytes:    1 << 20,
}
```

上传接口需要较长 `ReadTimeout`，但不能无限制。

### 7.7 修复缩略图更新错误判断

当前代码：

```go
_, thumbErr := fileService.fileDao.UpdateRowThumbPathByFileId(uploadDir.Thumb, fileId, userId)
if compErr != nil {
    log.Printf("%v: %v", Err6265, thumbErr)
}
```

这里应判断 `thumbErr`：

```go
if thumbErr != nil {
    log.Printf("%v: %v", Err6265, thumbErr)
}
```

这是可观测性问题，不是 OOM 主因，但建议一起修。

## 8. 建议的配置扩展

当前 `FileConfig` 只有路径和 host：

```go
type FileConfig struct {
    Path     string
    Host     string
    Original string
    Thumb    string
}
```

建议扩展为：

```go
type FileConfig struct {
    Path                 string `mapstructure:"path"`
    Host                 string `mapstructure:"host"`
    Original             string `mapstructure:"original"`
    Thumb                string `mapstructure:"thumb"`
    MaxFileSizeMB        int64  `mapstructure:"maxFileSizeMB"`
    MultipartMemoryMB    int64  `mapstructure:"multipartMemoryMB"`
    MaxUploadConcurrency int    `mapstructure:"maxUploadConcurrency"`
    ThumbWorkers         int    `mapstructure:"thumbWorkers"`
    ThumbQueueSize       int    `mapstructure:"thumbQueueSize"`
    MaxImagePixels       int    `mapstructure:"maxImagePixels"`
}
```

老笔记本推荐起始配置：

```yaml
file:
  thumb: "thumb"
  original: "original"
  path: "/var/www/bee_album"
  host: "http://192.168.2.3:6003"
  maxFileSizeMB: 10
  multipartMemoryMB: 1
  maxUploadConcurrency: 3
  thumbWorkers: 1
  thumbQueueSize: 50
  maxImagePixels: 20000000
```

如果实际测试内存仍高，优先把 `maxUploadConcurrency` 降到 2，不要增加 `thumbWorkers`。

## 9. 分阶段实施计划

### 阶段 1：止血，半天内完成

目标：上传 200 张图片时不再宕机。

改动：

1. 前端上传并发从 10 降到 2 或 3。
2. `ParseMultipartForm(32 << 20)` 改为 `ParseMultipartForm(1 << 20)`。
3. 请求体上限改为 `maxFileSize + 1MB`。
4. 缩略图生成加全局信号量，限制为 1 个并发。
5. 修复 `thumbErr` 日志判断。
6. 生产配置中 DB 连接池调低到 `maxOpenConns: 4`、`maxIdleConns: 2`。

预期效果：

- 上传时 Go 堆内请求体内存明显降低。
- 缩略图解码最多 1 张图同时进行。
- 200 张批量上传会变慢，但稳定性会显著提升。

### 阶段 2：可靠化，1 到 2 天

目标：后端能主动保护自己，不依赖前端克制。

改动：

1. 上传接口增加有界并发限制，过载返回 `429`。
2. 缩略图处理改成 worker pool + 有界队列。
3. 图片保存前后增加格式白名单、MIME 校验、`DecodeConfig` 像素上限校验。
4. `bee_file` 增加缩略图状态字段，例如 `thumb_status`：`pending`、`done`、`failed`。
5. 增加后台重试或手动修复缩略图任务。
6. HTTP Server 增加超时配置。

预期效果：

- 后端形成稳定背压。
- 无论前端开多少并发，服务端都不会无限创建图片处理任务。
- 缩略图失败可以被追踪和修复。

### 阶段 3：低内存优化，按需要实施

目标：在低配硬件上进一步提升吞吐和稳定性。

可选改动：

1. 上传解析改为 `MultipartReader` 流式落盘。
2. 缩略图生成改用更低内存的图像处理方案，例如 libvips/govips。
3. 将缩略图任务拆成独立进程，主服务只负责上传和 DB，图片处理进程单独设置内存限制。
4. 增加 `/debug/pprof` 或 Prometheus 指标，观察 heap、goroutine、缩略图队列长度、处理耗时。
5. 根据实际图片尺寸，把缩略图最长边降低到 512 或 720。

## 10. 验证方案

### 10.1 功能验证

覆盖场景：

1. 上传 1 张小于 10MB 的 jpg。
2. 上传 1 张小于 10MB 的 png。
3. 上传接近 10MB 的图片，确认不会因 multipart 开销被误拒。
4. 上传超过 10MB 的文件，返回 413。
5. 上传非图片文件，返回业务错误。
6. 上传超高像素但文件体积较小的图片，返回像素超限。
7. 上传成功后文件落盘、DB 插入、缩略图生成、`file_thumb_path` 更新正常。

### 10.2 压测验证

推荐压测条件：

```text
图片数量：200
单文件大小：1MB 到 10MB 混合
前端或压测并发：10
服务端上传并发限制：2 到 3
缩略图 worker：1
```

需要观察：

- Go 进程 RSS。
- 系统剩余内存和 swap 使用。
- goroutine 数量。
- 缩略图队列长度。
- 每张图上传耗时和缩略图耗时。
- 失败请求是否为可预期的 429、413 或格式错误。

通过标准：

```text
200 张图片上传期间服务不宕机。
Go 进程 RSS 长时间稳定在目标内存以内。
缩略图队列最终归零。
所有合法图片最终有原图记录，缩略图成功或有明确失败状态。
```

### 10.3 建议的运行参数

老笔记本初始建议：

```shell
GOMEMLIMIT=1200MiB GOGC=50 ./bee_go -env prod
```

如果机器还运行 MySQL、桌面环境或其他服务，不建议把 Go 服务内存目标设置太高。稳定后可以根据监控结果调整到 1500MiB 左右。

## 11. 推荐最终状态

对当前项目和硬件，推荐最终架构是：

```text
HTTP 上传接口
  -> 后端上传并发限制
  -> 低内存 multipart 解析或流式落盘
  -> 文件大小、格式、像素数校验
  -> 原图落盘
  -> DB 写入 thumb_status=pending
  -> 投递缩略图任务到有界队列

缩略图 worker
  -> 固定 1 个 worker
  -> 读取原图
  -> DecodeConfig 校验
  -> imaging/govips 生成缩略图
  -> 更新 thumb_status 和 file_thumb_path
```

关键目标不是让 4GB 机器“更快地同时处理更多图片”，而是让它在面对批量上传时有明确上限：上传可排队，缩略图可延迟，但服务不能因为内存失控而宕机。

## 12. 优先级排序

最推荐的实施顺序：

1. 限制缩略图生成并发。
2. 降低 multipart 内存阈值。
3. 增加图片像素上限校验。
4. 增加上传接口并发限制和 429 背压。
5. 设置 HTTP Server 超时。
6. 调低 DB 连接池。
7. 改为流式 multipart 上传。
8. 建立缩略图状态和可重试任务机制。

只做第 1、2、3 项，内存风险就会大幅下降；做到第 4 项后，后端才算具备基本自我保护能力。
