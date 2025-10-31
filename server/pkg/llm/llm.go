// pkg/llm/llm.go
package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

// RAGClient RAG 客户端
type RAGClient struct {
	llm            llms.Model
	promptTemplate string
	temperature    float64
}

// NewRAGClient 创建新的 RAG 客户端
func NewRAGClient(temperature float64) (*RAGClient, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, fmt.Errorf("初始化 LLM 失败: %w", err)
	}

	promptTemplate := `使用以下上下文回答问题。如果不知道，请说"无法回答"。

上下文：
{{.context}}

问题：{{.question}}

回答：`

	return &RAGClient{
		llm:            llm,
		promptTemplate: promptTemplate,
		temperature:    temperature,
	}, nil
}

// Answer 根据检索到的文档回答问题
func (r *RAGClient) Answer(ctx context.Context, question string, relevantDocs []schema.Document) (string, error) {
	// 构建上下文
	context := ""
	for i, doc := range relevantDocs {
		context += fmt.Sprintf("片段 %d: %s\n\n", i+1, doc.PageContent)
	}

	// 构造 prompt
	prompt := prompts.NewPromptTemplate(r.promptTemplate, []string{"context", "question"})
	promptVal, err := prompt.Format(map[string]any{
		"context":  context,
		"question": question,
	})
	if err != nil {
		return "", fmt.Errorf("格式化 prompt 失败: %w", err)
	}

	// 调用 LLM
	answer, err := r.llm.Call(ctx, promptVal, llms.WithTemperature(r.temperature))
	if err != nil {
		return "", fmt.Errorf("LLM 调用失败: %w", err)
	}

	return answer, nil
}
