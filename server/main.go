// main.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

// CustomEmbedder 实现 embeddings.Embedder
type CustomEmbedder struct {
	Endpoint string
}

func (e *CustomEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	return e.embed(texts)
}

func (e *CustomEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embs, err := e.embed([]string{text})
	if err != nil {
		return nil, err
	}
	return embs[0], nil
}

func (e *CustomEmbedder) embed(texts []string) ([][]float32, error) {
	reqBody, _ := json.Marshal(map[string][]string{"texts": texts})
	resp, err := http.Post(e.Endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call embed server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embed server returned status: %d", resp.StatusCode)
	}

	var result struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Embeddings, nil
}

var _ embeddings.Embedder = &CustomEmbedder{}

// loadMarkdownFiles 递归加载目录下所有 .md 文件
func loadMarkdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, string(content))
		}
		return nil
	})
	return files, err
}

func main() {
	ctx := context.Background()

	// 检查文档目录
	if _, err := os.Stat("../notion_docs"); os.IsNotExist(err) {
		log.Fatal("❌ 请将 Notion Markdown 文件放在 ./notion_docs 目录中")
	}

	// === 1. 加载所有 .md 文件内容 ===
	fmt.Println("📂 加载 Notion 文档...")
	rawContents, err := loadMarkdownFiles("../notion_docs")
	if err != nil {
		log.Fatal("加载文档失败:", err)
	}
	fmt.Printf("✅ 成功加载 %d 个文件\n", len(rawContents))

	// === 2. 使用 Markdown 分块器分块 ===
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(400),
		textsplitter.WithChunkOverlap(50),
		textsplitter.WithSeparators([]string{"\n\n", "\n", " ", ""}),
	)

	var allChunks []schema.Document
	for _, content := range rawContents {
		// 先用 Markdown splitter 初步处理（可选）
		mdSplitter := textsplitter.NewMarkdownTextSplitter()
		mdDocs, err := mdSplitter.SplitText(content)
		if err != nil {
			log.Printf("Warning: Markdown split failed, using raw text")
			mdDocs = []string{content}
		}

		// 再用递归分块
		for _, mdChunk := range mdDocs {
			chunks, err := splitter.SplitText(mdChunk)
			if err != nil {
				continue
			}
			for _, chunk := range chunks {
				if strings.TrimSpace(chunk) != "" {
					allChunks = append(allChunks, schema.Document{PageContent: chunk})
				}
			}
		}
	}
	fmt.Printf("✂️  分块完成，共 %d 个文本块\n", len(allChunks))

	// === 3. 创建 Embedder ===
	embedder := &CustomEmbedder{Endpoint: "http://localhost:8081/embed"}

	// === 4. 连接 Chroma ===
	fmt.Println("🔗 连接 Chroma 向量数据库...")
	store, err := chroma.New(
		chroma.WithNameSpace("notion_rag"),
		chroma.WithEmbedder(embedder),
		chroma.WithChromaURL("http://localhost:8000"),
	)
	if err != nil {
		log.Fatal("连接 Chroma 失败:", err)
	}

	// === 5. 导入文档（首次运行）===
	// 简单策略：尝试检索，失败则认为是首次
	fmt.Println("📥 检查是否需要导入文档...")
	_, err = store.SimilaritySearch(ctx, "hello", 1)
	if err != nil && (strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "collection")) {
		fmt.Println("🆕 首次运行，正在导入文档...")
		_, err := store.AddDocuments(ctx, allChunks)
		if err != nil {
			log.Printf("警告：导入文档时出现错误: %v", err)
		}
		fmt.Println("✅ 文档导入完成！")
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("✅ 文档已存在，跳过导入")
	}

	// === 6. 初始化 LLM ===
	// 需要设置环境变量 MOONSHOT_API_KEY
	fmt.Println("🧠 初始化 LLM (Kimi K2)...")
	apiKey := os.Getenv("MOONSHOT_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置环境变量 MOONSHOT_API_KEY")
	}
	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithBaseURL("https://api.moonshot.ai/v1"),
		openai.WithModel("kimi-k2-0711-preview"),
	)
	if err != nil {
		log.Fatal("初始化 LLM 失败:", err)
	}

	// === 7. 提问与检索 ===
	question := "如何高效进行项目复盘？"
	fmt.Printf("\n❓ 问题: %s\n", question)

	relevantDocs, err := store.SimilaritySearch(ctx, question, 3)
	if err != nil {
		log.Fatal("检索失败:", err)
	}

	context := ""
	for i, doc := range relevantDocs {
		context += fmt.Sprintf("片段 %d: %s\n\n", i+1, doc.PageContent)
	}

	// 构造 prompt
	promptTemplate := `使用以下上下文回答问题。如果不知道，请说“无法回答”。

上下文：
{{.context}}

问题：{{.question}}

回答：`
	prompt := prompts.NewPromptTemplate(promptTemplate, []string{"context", "question"})
	promptVal, _ := prompt.Format(map[string]any{
		"context":  context,
		"question": question,
	})

	// 调用 LLM（注意：WithTemperature 是 llms 包的选项）
	answer, err := llm.Call(ctx, promptVal, llms.WithTemperature(0.3)) // ✅ 修正
	if err != nil {
		log.Fatal("LLM 调用失败:", err)
	}

	fmt.Println("\n📚 检索结果:")
	fmt.Println(context)
	fmt.Println("🤖 回答:")
	fmt.Println(answer)
}
