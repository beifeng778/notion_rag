# Notion RAG 系统

<div align="center">

**基于 Go + LangChain + Chroma 的智能文档问答系统**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 🎯 项目简介

支持 Notion Markdown 文档的检索增强生成（RAG）系统，在 Windows 本地运行，16GB 内存友好，无需 GPU。

### ✨ 核心特性

- 🚀 **高性能**：Go 语言编写，内存占用低
- 🧠 **智能问答**：支持任何兼容 OpenAI API 的大语言模型
- 📚 **文档管理**：自动加载和分块 Notion 文档
- 🔍 **语义检索**：基于 Chroma 向量数据库的相似度搜索
- 💾 **本地优先**：除 LLM 外完全本地运行，数据安全

## 📁 项目结构

```text
notion_rag/
├── .venv/                  # Python 虚拟环境
├── chroma_db/              # Chroma 向量数据库持久化目录（自动生成）
├── notion_docs/            # Notion 文档目录
├── cmd/
│   └── embed_server.py     # Embedding 服务：BAAI/bge-m3 模型
├── server/
│   ├── cmd/chat/
│   │   └── main.go         # 终端对话主程序
│   ├── pkg/
│   │   ├── embedder/       # 向量化模块
│   │   ├── loader/         # 文档加载模块
│   │   ├── vectorstore/    # 向量数据库模块
│   │   └── llm/            # LLM 客户端模块
│   ├── env.sh              # 环境变量配置脚本
│   ├── README.md           # Server 端详细文档
│   ├── go.mod
│   └── go.sum
└── README.md               # 本文件
```

---

## ✅ 功能亮点

- ✅ 支持 **Notion Markdown 文档导入**（实际上任何 Markdown 文档都可以，其他的文档格式没有测试过）
- ✅ 使用 **BAAI/bge-m3** 作为 Embedding 模型（中文强、1024维）
- ✅ 使用 **Chroma** 作为向量数据库（本地持久化）
- ✅ 支持任何**兼容 OpenAI API 的大语言模型**（Moonshot Kimi、OpenAI GPT、DeepSeek 等）
- ✅ 完全本地运行（除 LLM 外），数据不出本地
- ✅ 16GB 内存机器友好，无需 GPU
- ✅ 终端对话模式，实时问答

---

## 🛠️ 技术栈与版本

### 核心依赖

| 组件 | 版本 | 说明 |
|------|------|------|
| **Go** | 1.24.5+ | 主程序语言 |
| **Python** | 3.13.3 | Embedding 服务 |
| **LangChain Go** | v0.1.14 | Go 语言的 LangChain 实现 |
| **Chroma** | 0.4.24 | 向量数据库 (v1 API) |
| **sentence-transformers** | 5.1.2 | Embedding 模型库 |
| **Flask** | 3.1.2 | Embedding 服务框架 |
| **PyTorch** | 2.9.0+cpu | 深度学习框架 (CPU版本) |
| **BAAI/bge-m3** | latest | 中文 Embedding 模型 (1024维) |
| **LLM** | 任意 | 支持任何兼容 OpenAI API 的模型 |

---

## 📦 安装教程

### 步骤 1: 安装 Go 环境

1. 下载并安装 [Go 1.24.5+](https://go.dev/dl/)
2. 验证安装:
   ```bash
   go version
   ```
   应该显示: `go version go1.24.5 windows/amd64`

### 步骤 2: 安装 Python 环境

1. 下载并安装 [Python 3.13+](https://www.python.org/downloads/windows/)
2. 验证安装:
   ```bash
   python --version
   ```
   应该显示: `Python 3.13.x` 或更高

### 步骤 3: 创建 Python 虚拟环境

在项目根目录 (`notion_rag/`) 执行:

```bash
# 创建虚拟环境
python -m venv .venv

# 激活虚拟环境
# Windows CMD:
.venv\Scripts\activate.bat

# Windows PowerShell:
.venv\Scripts\Activate.ps1

# Git Bash / Linux / Mac:
source .venv/Scripts/activate
```

### 步骤 4: 安装 Python 依赖

激活虚拟环境后,安装依赖:

```bash
pip install --upgrade pip
pip install sentence-transformers==5.1.2
pip install flask==3.1.2
pip install chromadb==0.4.24
```

### 步骤 5: 安装 Chroma 数据库

```bash
pip install chromadb==0.4.24
```

⚠️ **重要**: 必须使用 0.4.24 版本,因为 langchaingo 的 chroma 客户端只支持 v1 版本的 API。

验证安装:
```bash
chroma --version
```

### 步骤 6: 下载 Embedding 模型

首次运行 `embed_server.py` 时会自动下载 `BAAI/bge-m3` 模型 (~1.3GB),请耐心等待。

或手动预下载:
```python
from sentence_transformers import SentenceTransformer
model = SentenceTransformer('BAAI/bge-m3')
```

### 步骤 7: 安装 Go 依赖

进入 `server/` 目录:
```bash
cd server
go mod tidy
```

这会自动下载所有 Go 依赖,包括:
- `github.com/tmc/langchaingo` - LangChain Go 实现
- 其他依赖见 `go.mod`

---

## 🔧 配置环境变量

编辑 `server/env.sh`,填入你的配置:

```bash
#!/bin/bash

# LLM 配置 (支持任何兼容 OpenAI API 的模型)
export OPENAI_API_KEY="sk-your-api-key"
export OPENAI_MODEL="your-model-name"
export OPENAI_BASE_URL="https://your-api-base-url"

# 向量数据库配置
export CHROMA_URL="http://localhost:8000"

# Embedding 服务配置
export EMBED_ENDPOINT="http://localhost:8081/embed"

# 文档目录配置
export DOCS_DIR="../notion_docs"
```

---

## 🚀 运行步骤

> ⚠️ 需要打开 **3 个终端窗口**,按顺序启动服务!

### 终端 1: 启动 Chroma 向量数据库

```bash
cd notion_rag
chroma run --path ./chroma_db
```

✅ 成功输出:
```
INFO:     Uvicorn running on http://0.0.0.0:8000
```

### 终端 2: 启动 Embedding 服务

```bash
cd notion_rag
source .venv/Scripts/activate  # 激活虚拟环境
cd cmd
python embed_server.py
```

✅ 成功输出:
```
 * Running on http://127.0.0.1:8081
```

首次运行会下载 `BAAI/bge-m3` 模型 (~2GB),请耐心等待。

### 终端 3: 运行 Go 程序

```bash
cd notion_rag/server

# 加载环境变量
source env.sh

# 运行程序
cd cmd/chat && go run main.go
```

✅ 成功输出:
```
🚀 初始化 Notion RAG 系统...
📂 加载 Notion 文档...
✅ 成功加载并分块，共 108 个文本块
🔗 连接 Chroma 向量数据库...
📥 检查是否需要导入文档...
🆕 首次运行，正在导入文档...
✅ 文档导入完成！
🧠 初始化 LLM...

=============================================================
🎉 系统初始化完成！现在可以开始提问了
💡 输入 'exit' 或 'quit' 退出程序
=============================================================

❓ 你的问题:
```

---

## 📥 导入 Notion 文档

1. 在 Notion 中导出页面为 **Markdown & CSV**
2. 解压后，将所有 `.md` 文件放入：
   ```
   notion_rag/notion_docs/
   ```
3. 首次运行 Go 程序时会自动导入文档到 Chroma

> 💡 如果更新了文档，请删除 `chroma_db/` 目录，重新运行程序即可重新导入。

---

## 📊 性能与资源占用

| 组件 | 内存占用 | CPU | 是否需要 GPU |
|------|----------|-----|--------------|
| Chroma | ~100MB | 低 | ❌ |
| Embedding (bge-m3) | ~2GB | 中 | ❌ |
| LLM (API) | 0MB（远程） | 无 | ❌ |
| Go 主程序 | ~200MB | 低 | ❌ |
| **总计** | **~2.3GB** | 低 | ❌ |

> ✅ 16GB 内存机器完全无压力！

---

## 📝 注意事项

- **工作目录**：所有命令默认在项目根目录（`notion_rag/`）执行
- **首次运行较慢**：需要下载 Embedding 模型 (~2GB) 并嵌入所有文档
- **确保 `chroma_db/` 存在且可写**
- **LLM API 可能有调用次数限制**，请合理使用
- **如遇端口冲突**，可修改 `embed_server.py` 的端口和 `env.sh` 中的配置

---

## 🆘 常见问题

### Q: 报错 `410 Gone` 或 `connection refused`？

**A:** Chroma 数据库未启动或无法连接。

**解决方案：**
1. 确保 Chroma 服务已启动：
   ```bash
   chroma run --path ../chroma_db
   ```
2. 检查端口 8000 是否被占用：
   ```bash
   netstat -ano | findstr :8000
   ```
3. 测试 Chroma 连接：
   ```bash
   curl http://localhost:8000/api/v1/heartbeat
   ```

### Q: 报错 `请将 Notion Markdown 文件放在 ./notion_docs 目录中`？

**A:** 文档目录不存在或路径不正确。

**解决方案：**
1. 在项目根目录创建 `notion_docs` 文件夹
2. 将 Notion 导出的 `.md` 文件放入该目录
3. 确保从 `server/` 目录运行程序

### Q: Embedding 服务报错？

**A:** Python 依赖未安装或虚拟环境未激活。

**解决方案：**
```bash
source .venv/Scripts/activate  # 激活虚拟环境
pip install sentence-transformers flask chromadb
```

### Q: LLM 返回错误？

**A:** 检查 API Key 和模型配置。

**解决方案：**
1. 检查 `OPENAI_API_KEY` 环境变量是否正确设置
2. 验证 `OPENAI_BASE_URL` 和 `OPENAI_MODEL` 是否正确
3. 检查 API 配额是否充足

---

## 🚀 未来扩展建议

- 添加 Web UI（用 Gin/Fiber）
- 支持命令行输入问题
- 自动监控 Notion 更新并增量导入
- 集成 Redis 缓存
- 支持多用户、权限管理

---

## 🙏 致谢

感谢你使用本项目！如果你觉得有用，欢迎 Star ⭐ 或 Fork 🍴！

如有任何问题，欢迎提交 Issue 或联系我！

---

## 📜 许可证

MIT License

---

## 🔗 相关链接

- [Chroma 文档](https://docs.trychroma.com/)
- [LangChain Go](https://github.com/tmc/langchaingo)
- [Sentence Transformers](https://www.sbert.net/)
- [BAAI/bge-m3 模型](https://huggingface.co/BAAI/bge-m3)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star！⭐**

Made with ❤️ by rainyday

</div>