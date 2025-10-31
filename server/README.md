# Notion RAG 系统

基于 LangChain Go 和 Chroma 的 Notion 文档问答系统，支持终端对话模式。

## 项目结构

```
server/
├── cmd/
│   └── chat/          # 终端对话主程序
│       └── main.go
├── pkg/
│   ├── embedder/      # 向量化模块
│   │   └── embedder.go
│   ├── loader/        # 文档加载模块
│   │   └── loader.go
│   ├── vectorstore/   # 向量数据库模块
│   │   └── vectorstore.go
│   └── llm/          # LLM 客户端模块
│       └── llm.go
├── env.sh            # 环境变量配置脚本
├── go.mod
└── go.sum
```

## 系统工作流程

### 📥 初始化阶段（首次运行）

```
1. loader 模块
   ↓ 加载 Markdown 文档
   ↓ 文档分块处理（支持递归分块和 Markdown 分块）
   ↓ 输出: []schema.Document（文档块列表）

2. embedder 模块
   ↓ 调用本地 embedding 服务（http://localhost:8081/embed）
   ↓ 使用 BAAI/bge-m3 模型向量化文档块
   ↓ 输出: 1024 维向量

3. vectorstore 模块
   ↓ 调用本地 Chroma 向量数据库（http://localhost:8000）
   ↓ 存储文档块和对应的向量
   ↓ 建立向量索引
```

### 💬 问答阶段（用户提问）

```
1. 用户输入问题
   ↓ 例如: "如何进行项目复盘？"

2. vectorstore 模块
   ↓ 调用 store.SimilaritySearch(ctx, query, topK)
   ↓ LangChain 内部自动调用 embedder 将问题向量化
   ↓ 在 Chroma 中进行相似度搜索（余弦相似度）
   ↓ 返回 Top-K 个最相关的文档块
   ↓ 输出: []schema.Document（相关文档列表）

3. llm 模块
   ↓ 构建 Prompt（包含用户问题 + 检索到的上下文）
   ↓ 调用 LLM API（支持任何兼容 OpenAI API 的模型）
   ↓ 输出: 基于上下文的答案
```

## 模块详细说明

### embedder (向量化模块)
- **职责**: 负责文本向量化
- **实现**: 调用本地 embedding 服务（Flask + sentence-transformers）
- **模型**: BAAI/bge-m3（1024维，中文强）
- **接口**: 实现 LangChain 的 Embedder 接口

### loader (文档加载模块)
- **职责**: 加载和分块 Markdown 文档
- **支持格式**: `.md` 文件
- **分块策略**: 
  - Markdown 分块（按标题、段落）
  - 递归分块（按字符数，支持重叠）
- **输出**: `[]schema.Document`

### vectorstore (向量数据库模块)
- **职责**: 封装 Chroma 向量数据库操作
- **功能**:
  - 文档初始化和导入
  - 相似度搜索（余弦相似度）
  - 强制重新导入（`FORCE_REIMPORT`）
- **存储**: 本地持久化到 `chroma_db/` 目录

### llm (LLM 客户端模块)
- **职责**: 封装 LLM API 调用和 RAG 问答逻辑
- **支持模型**: 任何兼容 OpenAI API 的模型
  - Moonshot Kimi
  - OpenAI GPT
  - DeepSeek
  - 其他兼容模型
- **功能**:
  - Prompt 模板管理
  - 上下文注入
  - 流式/非流式响应

## 环境配置

本项目使用环境变量脚本管理配置,支持灵活切换不同的 LLM 模型。

### 创建 env.sh 文件

在 `server` 目录下创建 `env.sh`:
```bash
#!/bin/bash
# 环境变量配置脚本
# 使用方式: source env.sh

# ============================================
# LLM 配置 (支持多种模型)
# ============================================

# Moonshot Kimi (默认)
export OPENAI_API_KEY="sk-your-api-key-here"
export OPENAI_MODEL="kimi-k2-0711-preview"
export OPENAI_BASE_URL="https://api.moonshot.cn/v1"

# 如果要切换到其他模型,注释掉上面的配置,取消注释下面对应的配置

# OpenAI GPT-4
# export OPENAI_API_KEY="sk-your-openai-key"
# export OPENAI_MODEL="gpt-4"
# export OPENAI_BASE_URL="https://api.openai.com/v1"

# DeepSeek
# export OPENAI_API_KEY="sk-your-deepseek-key"
# export OPENAI_MODEL="deepseek-chat"
# export OPENAI_BASE_URL="https://api.deepseek.com/v1"

# ============================================
# 向量数据库配置
# ============================================
export CHROMA_URL="http://localhost:8000"

# ============================================
# Embedding 服务配置
# ============================================
export EMBED_ENDPOINT="http://localhost:8081/embed"

# ============================================
# 文档目录配置
# ============================================
export DOCS_DIR="../notion_docs"

echo "✅ 环境变量已加载"
echo "📌 当前配置:"
echo "   - Model: $OPENAI_MODEL"
echo "   - Base URL: $OPENAI_BASE_URL"
echo "   - Chroma: $CHROMA_URL"
echo "   - Embed: $EMBED_ENDPOINT"
echo "   - Docs: $DOCS_DIR"
```

### 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `OPENAI_API_KEY` | LLM API 密钥 | 必填 |
| `OPENAI_MODEL` | 模型名称 | `kimi-k2-0711-preview` |
| `OPENAI_BASE_URL` | API 基础 URL | `https://api.moonshot.cn/v1` |
| `CHROMA_URL` | Chroma 数据库地址 | `http://localhost:8000` |
| `EMBED_ENDPOINT` | Embedding 服务地址 | `http://localhost:8081/embed` |
| `DOCS_DIR` | 文档目录路径 | `../notion_docs` |

## 运行方式

### 1. 配置环境变量

创建并编辑 `env.sh`,填入你的 API Key (参考上面的完整示例)

### 2. 准备文档

将 Notion 导出的 Markdown 文档放到 `notion_docs` 目录:
```
notion_rag/
├── notion_docs/      # 在这里放置你的 .md 文件
│   ├── doc1.md
│   ├── doc2.md
│   └── ...
└── server/
```

### 3. 启动必要服务

确保以下服务已启动:
- Chroma 向量数据库 (默认端口 8000)
- Embedding 服务 (默认端口 8081)

### 4. 运行程序

```bash
# 1. 加载环境变量
source env.sh

# 2. 运行程序
cd cmd/chat && go run main.go
```

⚠️ **重要**: 必须先执行 `source env.sh` 加载环境变量,否则程序会报错找不到 API Key

## 使用说明

1. 程序启动后会自动加载 `DOCS_DIR` 指定目录下的所有 Markdown 文件 (默认 `../notion_docs`)
2. 首次运行会将文档导入向量数据库
3. 进入对话模式后，直接输入问题即可
4. 输入 `exit` 或 `quit` 退出程序

## 切换 LLM 模型

编辑 `env.sh`,注释当前模型配置,取消注释目标模型:

```bash
# 注释掉 Moonshot
# export OPENAI_API_KEY="sk-moonshot-key"
# export OPENAI_MODEL="kimi-k2-0711-preview"
# export OPENAI_BASE_URL="https://api.moonshot.cn/v1"

# 使用 OpenAI GPT-4
export OPENAI_API_KEY="sk-openai-key"
export OPENAI_MODEL="gpt-4"
export OPENAI_BASE_URL="https://api.openai.com/v1"
```

然后重新加载环境变量:
```bash
source env.sh
```

## 示例对话

```
🚀 初始化 Notion RAG 系统...

📂 加载 Notion 文档...
✅ 成功加载并分块，共 108 个文本块
🔗 连接 Chroma 向量数据库...
📥 检查是否需要导入文档...
✅ 文档已存在，跳过导入
🧠 初始化 LLM (Kimi K2)...

============================================================
🎉 系统初始化完成！现在可以开始提问了
💡 输入 'exit' 或 'quit' 退出程序
============================================================

❓ 你的问题: 如何高效进行项目复盘？

🔍 正在检索相关文档...
📚 找到 3 个相关片段

🤖 正在生成回答...

💬 回答:
------------------------------------------------------------
根据提供的上下文，高效进行项目复盘的关键步骤包括...
------------------------------------------------------------

❓ 你的问题: exit
👋 再见！
```

## 技术栈

- **Go**: 1.24.5
- **LangChain Go**: 文档处理和 LLM 集成
- **Chroma**: 向量数据库
- **LLM**: 支持 Moonshot Kimi, OpenAI GPT-4, DeepSeek 等兼容 OpenAI API 的模型

## 常见问题

### Q: 提示找不到 OPENAI_API_KEY?
**A**: 确保运行程序前先执行了 `source env.sh` 加载环境变量

### Q: 如何修改文档目录?
**A**: 在 `env.sh` 中修改 `DOCS_DIR` 变量:
```bash
export DOCS_DIR="/path/to/your/docs"
```

### Q: 支持哪些 LLM 模型?
**A**: 支持所有兼容 OpenAI API 格式的模型,只需配置对应的 API Key、Model 和 Base URL 即可

### Q: 中文检索不到结果怎么办?
**A**: 这通常是因为 Chroma 中的旧数据使用了不同的 embedding 模型。解决方法:

1. 编辑 `env.sh`,取消注释 `FORCE_REIMPORT`:
```bash
export FORCE_REIMPORT="true"
```

2. 重新加载环境变量并运行:
```bash
source env.sh
cd cmd/chat && go run main.go
```

3. 导入完成后,注释掉 `FORCE_REIMPORT` 避免每次都重新导入:
```bash
# export FORCE_REIMPORT="true"
```
