# 中文检索问题诊断指南

## 问题现象
查询中文内容时,Chroma 返回 0 个结果

## 可能原因
1. **Chroma 中的旧向量数据** - 使用了不同的 embedding 模型
2. **Embedding 服务问题** - 中文向量化失败
3. **文档分块问题** - 中文被错误分割
4. **Chroma collection 名称冲突** - 多个 collection 混用

## 诊断步骤

### 步骤 1: 测试 Embedding 服务

```bash
# 测试中文向量化
curl -X POST http://localhost:8081/embed \
  -H "Content-Type: application/json" \
  -d '{"texts": ["新版宜搭同步方案"]}' | python -m json.tool | head -20
```

**预期结果**: 应该看到一个包含 1024 维向量的数组

### 步骤 2: 检查 Chroma 数据

```bash
# 查看 Chroma 的 collections
curl http://localhost:8000/api/v1/collections | python -m json.tool
```

**预期结果**: 应该看到 `notion_rag` collection

### 步骤 3: 强制重新导入文档

1. 编辑 `env.sh`:
```bash
export FORCE_REIMPORT="true"
```

2. 运行程序:
```bash
source env.sh
cd cmd/chat && go run main.go
```

3. 观察输出:
- 应该看到 "🔄 强制重新导入模式..."
- 应该看到 "📝 正在导入 108 个文档块..."
- 应该看到每个文件的分块数量

### 步骤 4: 测试检索

重新导入后,查询 "宜搭",观察调试输出:
```
🔍 [DEBUG] 检索查询: '宜搭', TopK: 3
📊 [DEBUG] 检索到 X 个结果
   [1] 相似度分数: 0.xxxx, 内容预览: ...
```

## 常见问题

### Q: FORCE_REIMPORT 后还是 0 个结果?

**可能原因**: Chroma collection 没有被清空

**解决方案**:
1. 停止 Chroma
2. 删除 Chroma 数据目录
3. 重新启动 Chroma
4. 再次运行 FORCE_REIMPORT

### Q: Embedding 服务返回错误?

**检查**:
- embedding 服务是否在运行: `curl http://localhost:8081/embed`
- 模型是否正确加载: 查看 embed_server.py 的启动日志
- 端口是否正确: 检查 env.sh 中的 EMBED_ENDPOINT

### Q: 文档分块数量异常?

**检查**:
- 文档目录是否正确: `echo $DOCS_DIR`
- 文档是否存在: `ls -la $DOCS_DIR/*.md`
- 文档编码是否正确: 应该是 UTF-8

## 完整重置流程

如果以上都不行,执行完整重置:

```bash
# 1. 停止所有服务
# Ctrl+C 停止 Go 程序
# Ctrl+C 停止 Chroma
# Ctrl+C 停止 embed_server

# 2. 清理 Chroma 数据
# 找到 Chroma 数据目录并删除 (通常在 ~/.chroma 或当前目录)
rm -rf ~/.chroma

# 3. 重新启动服务
# 启动 Chroma
docker run -p 8000:8000 chromadb/chroma

# 启动 embedding 服务
cd cmd && python embed_server.py

# 4. 强制重新导入
cd server
export FORCE_REIMPORT="true"
source env.sh
cd cmd/chat && go run main.go

# 5. 测试
# 查询: "宜搭"
# 应该能看到结果
```

## 调试输出说明

程序现在会输出详细的调试信息:

```
📄 [DEBUG] 文件 PMP.md: 108 个文本块
🔍 [DEBUG] 检索查询: '宜搭', TopK: 3
📊 [DEBUG] 检索到 3 个结果
   [1] 相似度分数: 0.8234, 内容预览: 新版宜搭同步方案...
```

如果看不到这些输出,说明代码没有正确更新。
