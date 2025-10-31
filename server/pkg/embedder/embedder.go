// pkg/embedder/embedder.go
package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/embeddings"
)

// CustomEmbedder 实现 embeddings.Embedder 接口
type CustomEmbedder struct {
	Endpoint string
}

// NewCustomEmbedder 创建新的 CustomEmbedder 实例
func NewCustomEmbedder(endpoint string) *CustomEmbedder {
	return &CustomEmbedder{
		Endpoint: endpoint,
	}
}

// EmbedDocuments 对多个文档进行向量化
func (e *CustomEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	return e.embed(texts)
}

// EmbedQuery 对单个查询进行向量化
func (e *CustomEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embs, err := e.embed([]string{text})
	if err != nil {
		return nil, err
	}
	return embs[0], nil
}

// embed 内部方法,调用 embedding 服务
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

// 确保实现了 embeddings.Embedder 接口
var _ embeddings.Embedder = &CustomEmbedder{}
