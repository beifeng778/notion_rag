# Notion RAG 系统

<div align="center">

**基于 Go + LangChain + Chroma + Kimi K2 的智能文档问答系统**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 🎯 项目简介

支持 Notion Markdown 文档的检索增强生成（RAG）系统，在 Windows 本地运行，16GB 内存友好，无需 GPU。

### ✨ 核心特性

- 🚀 **高性能**：Go 语言编写，内存占用低
- 🧠 **智能问答**：集成 Kimi K2 大语言模型
- 📚 **文档管理**：自动加载和分块 Notion 文档
- 🔍 **语义检索**：基于 Chroma 向量数据库的相似度搜索
- 💾 **本地优先**：除 LLM 外完全本地运行，数据安全

## 📁 项目结构

```text
notion_rag/
├── .venv/                  # Python 虚拟环境
├── chroma_db/              # Chroma 向量数据库持久化目录（自动生成）
├── cmd/
│   └── embed_server.py     # 嵌入服务：all-MiniLM-L6-v2 模型
├── server/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go             # Go 主程序（RAG 核心逻辑）
│   └── notion_rag.exe      # 编译后的可执行文件（Windows）
└── README.md               # 本文件
```

---

## ✅ 功能亮点

- ✅ 支持 **Notion Markdown 文档导入**
- ✅ 使用 **all-MiniLM-L6-v2** 作为嵌入模型（轻量、英文好）
- ✅ 使用 **Chroma** 作为向量数据库（本地持久化）
- ✅ 使用 **Moonshot AI Kimi K2** 作为大语言模型（中文强、API 调用）
- ✅ 完全本地运行（除 LLM 外），数据不出本地
- ✅ 16GB 内存机器友好，无需 GPU

---

## 🛠️ 运行前准备

### 1. 安装依赖

#### Go 环境

- 安装 [Go 1.21+](https://go.dev/dl/)
- 验证：

  ```bash
  go version
  ```

#### Python 环境

- 安装 [Python 3.9+](https://www.python.org/downloads/windows/)
- 在项目根目录创建虚拟环境：

  ```bash
  python -m venv .venv
  .venv\Scripts\activate  # Windows
  # 或 source .venv/Scripts/activate  # Git Bash
  ```

#### 安装 Python 依赖

```bash
pip install sentence-transformers flask chromadb
```

#### 获取 Moonshot API Key

- 注册 [Moonshot AI 平台](https://platform.moonshot.ai/)
- 创建 API Key，保存为环境变量

---

## 🔧 设置环境变量（Windows）

在终端中设置（每次新终端需重新设置）：
```bash
export MOONSHOT_API_KEY="sk-xxx-your-key-here"
```

或永久设置（推荐）：
1. 按 `Win + R` → 输入 `sysdm.cpl` → 打开“系统属性”
2. “高级” → “环境变量”
3. 用户变量 → 新建 → 变量名：`MOONSHOT_API_KEY`，值：你的密钥
4. 重启终端生效

---

## 🚀 运行步骤（三步走）

> ⚠️ 请在项目根目录（`notion_rag/`）下打开三个终端窗口，按顺序运行以下命令！

### 🔹 终端 1：启动 Chroma 向量数据库
```bash
chroma run --path ./chroma_db
```

✅ 输出应包含：
```
INFO:     Uvicorn running on http://0.0.0.0:8000
```

### 🔹 终端 2：启动 Embedding 服务
```bash
cd cmd
python embed_server.py
```

✅ 输出应包含：
```
 * Running on http://127.0.0.1:8081
```

### 🔹 终端 3：运行 Go RAG 程序
```bash
cd server
set MOONSHOT_API_KEY=your_api_key_here
go run main.go
```

> 或先编译再运行：
> ```bash
> go build -o notion_rag
> ./notion_rag
> ```

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

## 🧪 示例问题

程序默认提问：
```go
question := "如何高效进行项目复盘？"
```

你可以修改 `main.go` 中的 `question` 变量来测试不同问题。

---

## 🔄 更新模型（可选）

### 换成更强中文模型：BAAI/bge-m3

编辑 `cmd/embed_server.py`：
```python
# 修改这一行
model = SentenceTransformer('BAAI/bge-m3')  # 中文更强，约 1.3GB
```

然后重新运行 `embed_server.py` 即可。

---

## 📊 性能与资源占用

| 组件 | 内存占用 | CPU | 是否需要 GPU |
|------|----------|-----|--------------|
| Chroma | ~100MB | 低 | ❌ |
| Embedding (all-MiniLM-L6-v2) | ~500MB | 中 | ❌ |
| Kimi K2 (API) | 0MB（远程） | 无 | ❌ |
| Go 主程序 | ~200MB | 低 | ❌ |
| **总计** | **< 1GB** | 低 | ❌ |

> ✅ 16GB 内存机器完全无压力！

---

## 📝 注意事项

- **工作目录**：所有命令默认在项目根目录（`notion_rag/`）执行
- **首次运行较慢**：需要下载 Embedding 模型并嵌入所有文档
- **确保 `chroma_db/` 存在且可写**
- **Moonshot API 有调用次数限制**，请合理使用
- **如遇端口冲突**，可修改 `embed_server.py` 的端口（如 8082）和 Go 代码中的 `Endpoint`

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

### Q: Kimi 返回乱码或错误？

**A:** API Key 配置问题。

**解决方案：**
1. 检查 `MOONSHOT_API_KEY` 环境变量是否正确设置
2. 在 [Moonshot 控制台](https://platform.moonshot.ai/) 验证 API Key 有效性
3. 检查 API 配额是否充足

### Q: 如何支持更多文件格式？

**A:** 当前只支持 `.md`，如需支持 `.txt` 或 `.pdf`，可扩展 `loadMarkdownFiles` 函数。

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

- [Moonshot AI 平台](https://platform.moonshot.ai/)
- [Chroma 文档](https://docs.trychroma.com/)
- [LangChain Go](https://github.com/tmc/langchaingo)
- [Sentence Transformers](https://www.sbert.net/)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star！⭐**

Made with ❤️ by [Your Name]

</div>