# claude-switch

`claude-switch` 是一个用于切换 Claude Code 后端供应商的 Go 命令行工具，安装后的命令名是 `cs`。

它会更新 Claude Code 使用的 `settings.json`，让你在不同兼容供应商之间快速切换，同时保留无关配置。

仓库地址：

```bash
git clone git@github.com:doublepi123/claude_switch.git
cd claude_switch
```

## 功能

- 支持列出当前内置供应商
- 支持查看当前 Claude Code 指向的供应商
- 支持保存各供应商 API Key
- 支持一键切换 `~/.claude/settings.json`
- 切换前自动备份原始配置
- 仅更新受管理的 `env` 字段，避免覆盖其他自定义配置

## 当前支持的供应商

| Provider | Base URL | 默认模型 |
| --- | --- | --- |
| `minimax` | `https://api.minimaxi.com/anthropic` | `MiniMax-M2.7` |
| `openrouter` | `https://openrouter.ai/api` | `anthropic/claude-sonnet-4.6` |
| `opencode-go` | `https://opencode.ai/zen/go` | `opencode-go/minimax-m2.7` |

## 安装

要求：

- Go 1.20 或更高版本

推荐先克隆仓库，再执行安装脚本。

### macOS / Linux

```bash
chmod +x scripts/install.sh
./scripts/install.sh
```

默认会安装到：

```text
~/.local/bin/cs
```

安装完成后可验证：

```bash
cs list
```

如需自定义安装目录：

```bash
INSTALL_DIR=/usr/local/bin ./scripts/install.sh
```

### Windows PowerShell

```powershell
.\scripts\install.ps1
```

默认会安装到：

```text
$HOME\AppData\Local\Programs\claude-switch\bin\cs.exe
```

安装完成后可验证：

```powershell
cs.exe list
```

如需自定义安装目录：

```powershell
.\scripts\install.ps1 -InstallDir 'C:\Tools\claude-switch'
```

### 手动构建

如果你不想使用安装脚本，也可以直接构建：

```bash
go build -o cs .
```

如果只想临时运行，也可以直接：

```bash
go run .
```

## 首次使用

第一次使用建议按下面顺序执行：

### 1. 查看支持的供应商

```bash
cs list
```

### 2. 保存 API Key

例如保存 MiniMax 的 Key：

```bash
cs set-key minimax sk-xxx
```

也可以保存 OpenRouter：

```bash
cs set-key openrouter sk-or-xxx
```

### 3. 切换到目标供应商

```bash
cs switch minimax
```

### 4. 确认当前生效配置

```bash
cs current
```

如果输出里能看到目标 `provider`、`base_url` 和 `model`，说明首次配置已完成。

## 命令说明

### 1. 查看可用供应商

```bash
cs list
```

输出包含供应商名称、Base URL 和默认模型。

### 2. 查看当前配置

```bash
cs current
```

默认读取：

```text
~/.claude/settings.json
```

也可以通过 `--claude-dir` 指定 Claude 配置目录：

```bash
cs current --claude-dir /path/to/.claude
```

### 3. 保存 API Key

```bash
cs set-key minimax sk-xxx
cs set-key openrouter sk-or-xxx
```

保存后会写入：

```text
~/.claude-switch/config.json
```

如果你不想落盘保存，也可以在切换时临时传入 `--api-key`。

### 4. 切换供应商

```bash
cs switch minimax
cs switch openrouter
cs switch opencode-go
```

如果没有提前保存 API Key，也可以在切换时直接传入：

```bash
cs switch openrouter --api-key sk-or-xxx
```

如果本地还没有 `~/.claude/settings.json`，工具会自动创建它。

### 5. 覆盖默认模型

```bash
cs switch opencode-go --model opencode-go/glm-5
```

### 6. 指定 Claude 配置目录

```bash
cs switch minimax --claude-dir /path/to/.claude
```

这对测试环境、多套 Claude 配置，或者首次调试很有用。

## 配置文件行为

默认情况下，工具会操作以下文件：

- Claude 配置：`~/.claude/settings.json`
- 本工具配置：`~/.claude-switch/config.json`

在执行 `switch` 时：

- 如果 `settings.json` 已存在，会先创建一个带时间戳的备份文件
- 只会清理并重写本工具管理的环境变量
- 其他字段和未受管理的环境变量会保留

当前受管理的环境变量包括：

- `ANTHROPIC_BASE_URL`
- `ANTHROPIC_AUTH_TOKEN`
- `ANTHROPIC_MODEL`
- `ANTHROPIC_DEFAULT_HAIKU_MODEL`
- `ANTHROPIC_DEFAULT_SONNET_MODEL`
- `ANTHROPIC_DEFAULT_OPUS_MODEL`
- `API_TIMEOUT_MS`
- `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC`

## 示例

先保存 Key：

```bash
cs set-key minimax sk-xxx
```

再切换：

```bash
cs switch minimax
cs current
```

直接临时切换而不保存 Key：

```bash
cs switch openrouter --api-key sk-or-xxx
```

切换并覆盖模型：

```bash
cs switch opencode-go --api-key your-key --model opencode-go/glm-5
```

## 测试

运行单元测试：

```bash
go test ./...
```

## 常见问题

### 1. 安装后找不到 `cs`

说明安装目录还没加入 `PATH`。

macOS / Linux 通常需要把下面路径加入 shell 配置：

```text
~/.local/bin
```

Windows 通常需要把下面路径加入用户 `Path`：

```text
$HOME\AppData\Local\Programs\claude-switch\bin
```

### 2. 切换前会不会覆盖我原来的 Claude 配置

不会整体覆盖。工具只会更新它自己管理的供应商相关环境变量，并在写入前自动备份已有的 `settings.json`。

### 3. 必须先执行 `set-key` 吗

不是必须。你也可以在 `switch` 时通过 `--api-key` 临时传入。

## 适用场景

- 在不同 Claude 兼容后端之间快速切换
- 为本地 Claude Code 环境维护统一的供应商配置
- 减少手动编辑 `settings.json` 的出错概率
