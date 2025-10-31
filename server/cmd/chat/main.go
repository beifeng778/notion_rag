// cmd/chat/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"notion_rag/pkg/embedder"
	"notion_rag/pkg/llm"
	"notion_rag/pkg/loader"
	"notion_rag/pkg/vectorstore"
)

func main() {
	ctx := context.Background()

	// 从环境变量读取文档目录
	docsDir := os.Getenv("DOCS_DIR")
	if docsDir == "" {
		log.Fatalf("❌ 未设置文档目录,请设置DOCS_DIR环境变量")
	}

	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		log.Fatalf("❌ 文档目录不存在: %s", docsDir)
	}

	fmt.Println("🚀 初始化 Notion RAG 系统...")
	fmt.Println()

	// === 1. 加载并分块文档 ===
	fmt.Println("📂 加载 Notion 文档...")
	docLoader := loader.NewDocumentLoader(docsDir, 400, 50)
	allChunks, err := docLoader.LoadAndSplit()
	if err != nil {
		log.Fatal("加载文档失败:", err)
	}
	fmt.Printf("✅ 成功加载并分块，共 %d 个文本块\n", len(allChunks))

	// === 2. 创建 Embedder ===
	embedEndpoint := os.Getenv("EMBED_ENDPOINT")
	if embedEndpoint == "" {
		embedEndpoint = "http://localhost:8081/embed"
	}
	embed := embedder.NewCustomEmbedder(embedEndpoint)

	// === 3. 连接向量数据库 ===
	fmt.Println("🔗 连接 Chroma 向量数据库...")
	chromaURL := os.Getenv("CHROMA_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	store, err := vectorstore.NewVectorStore(ctx, chromaURL, "notion_rag", embed)
	if err != nil {
		log.Fatal(err)
	}

	// === 4. 初始化文档（首次运行）===
	fmt.Println("📥 检查是否需要导入文档...")
	
	// 检查是否需要强制重新导入
	forceReimport := os.Getenv("FORCE_REIMPORT")
	if forceReimport == "true" {
		fmt.Println("🔄 强制重新导入模式...")
		if err := store.ForceReimportDocuments(ctx, allChunks); err != nil {
			log.Fatalf("❌ 强制重新导入失败: %v", err)
		}
	} else {
		if err := store.InitializeDocuments(ctx, allChunks); err != nil {
			log.Printf("警告：导入文档时出现错误: %v", err)
		}
	}

	// === 5. 初始化 LLM ===
	fmt.Println("🧠 初始化 LLM (Kimi K2)...")
	ragClient, err := llm.NewRAGClient(0.3)
	if err != nil {
		log.Fatal(err)
	}

	// === 6. 开始终端对话 ===
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("🎉 系统初始化完成！现在可以开始提问了")
	fmt.Println("💡 输入 'exit' 或 'quit' 退出程序")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("❓ 你的问题: ")

		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())

		// 检查退出命令
		if question == "exit" || question == "quit" {
			fmt.Println("👋 再见！")
			break
		}

		// 跳过空输入
		if question == "" {
			continue
		}

		fmt.Println()
		fmt.Println("🔍 正在检索相关文档...")

		// 检索相关文档
		relevantDocs, err := store.Search(ctx, question, 3)
		if err != nil {
			log.Printf("❌ 检索失败: %v\n", err)
			continue
		}

		fmt.Printf("📚 找到 %d 个相关片段\n", len(relevantDocs))
		fmt.Println()

		// 生成回答
		fmt.Println("🤖 正在生成回答...")
		answer, err := ragClient.Answer(ctx, question, relevantDocs)
		if err != nil {
			log.Printf("❌ 生成回答失败: %v\n", err)
			continue
		}

		fmt.Println()
		fmt.Println("💬 回答:")
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println(answer)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取输入时出错: %v\n", err)
	}
}
