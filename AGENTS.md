# AGENTS.md

## 项目概览

- 这是 github.com/harluo/taskd Go 模块，目标是提供基于数据库的任务执行框架。
- Go 版本以 go.mod 中声明的 go 1.25.0 为准。
- 根目录的 agent.go 只暴露 internal/core.Agent 别名；核心实现位于 internal/ 下，目录层级较深是现有项目结构，不要随意重排。
- internal/core/di.go 负责通过 github.com/harluo/di 注册依赖；数据库、迁移、服务和任务执行相关代码分别位于对应的 internal/internal/ 子目录。

## 开发约定

- 修改前先确认当前仓库根目录：git rev-parse --show-toplevel。
- 保留工作区中与当前任务无关的修改，不要回退、格式化或重写无关文件。
- 遵循现有包边界和依赖注入方式；优先复用现有类型与依赖，不为单一实现新增抽象层。
- 业务错误应沿用项目当前依赖和调用链的处理方式；不要为了局部改动引入新的错误、日志或配置框架。
- Go 代码使用 gofmt；修改导入后同步整理 import。
- 不要提交生成物、临时文件、密钥或本地配置。
- 注释和新增文档优先使用中文；代码中的公开 API 注释应清楚说明用途。

## 常用命令

所有外部命令按仓库约定使用 rtk 前缀：

    rtk go test ./...
    rtk go vet ./...
    rtk gofmt -w <修改过的.go文件>
    rtk git diff --check
    rtk git status --short

当前仓库没有 *_test.go 文件；行为发生变化时，优先在受影响包中补充最小必要测试，并至少运行 rtk go test ./...。

## 修改与验证流程

1. 先查看相关包、调用方和数据库/迁移影响范围，再进行最小范围修改。
2. 修改后检查 git diff，确认没有覆盖已有工作或产生无关格式化。
3. 对 Go 文件运行 gofmt，然后运行 rtk go test ./...；必要时再运行 rtk go vet ./...
4. 最后运行 rtk git diff --check 和 rtk git status --short，在回复中只报告实际执行过的检查。

## 目录提示

- agent.go：对外导出的任务代理别名。
- internal/core/：核心公开别名和依赖注入入口。
- internal/internal/model/：任务、调度和运行时模型。
- internal/internal/db/：数据库访问及内部列、SQL、查询辅助代码。
- internal/internal/migrate/：数据库迁移及表定义。
- internal/internal/service/：服务层。
- internal/internal/kernel/、internal/internal/put/、internal/internal/internal/：任务执行和内部适配逻辑。
