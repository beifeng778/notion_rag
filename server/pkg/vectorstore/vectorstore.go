// pkg/vectorstore/vectorstore.go
package vectorstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

// VectorStore 向量数据库封装
type VectorStore struct {
	store     chroma.Store
	namespace string
}

// NewVectorStore 创建新的向量数据库实例
func NewVectorStore(ctx context.Context, chromaURL, namespace string, embedder embeddings.Embedder) (*VectorStore, error) {
	store, err := chroma.New(
		chroma.WithNameSpace(namespace),
		chroma.WithEmbedder(embedder),
		chroma.WithChromaURL(chromaURL),
	)
	if err != nil {
		return nil, fmt.Errorf("连接 Chroma 失败: %w", err)
	}

	return &VectorStore{
		store:     store,
		namespace: namespace,
	}, nil
}

// InitializeDocuments 初始化文档（首次导入）
func (v *VectorStore) InitializeDocuments(ctx context.Context, docs []schema.Document) error {
	// 检查是否已有数据
	_, err := v.store.SimilaritySearch(ctx, "hello", 1)
	if err != nil && (strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "collection")) {
		fmt.Println("🆕 首次运行，正在导入文档...")
		_, err := v.store.AddDocuments(ctx, docs)
		if err != nil {
			return fmt.Errorf("导入文档失败: %w", err)
		}
		fmt.Println("✅ 文档导入完成！")
		time.Sleep(2 * time.Second)
		return nil
	}
	fmt.Println("✅ 文档已存在，跳过导入")
	return nil
}

// ForceReimportDocuments 强制重新导入文档（删除旧数据）
func (v *VectorStore) ForceReimportDocuments(ctx context.Context, docs []schema.Document) error {
	fmt.Println("🗑️  删除旧数据...")
	// Chroma 会自动处理 collection 的重建
	// 直接添加新文档即可，langchaingo 的 chroma 实现会处理
	
	fmt.Printf("📝 正在导入 %d 个文档块...\n", len(docs))
	_, err := v.store.AddDocuments(ctx, docs)
	if err != nil {
		return fmt.Errorf("导入文档失败: %w", err)
	}
	
	fmt.Println("✅ 文档重新导入完成！")
	time.Sleep(2 * time.Second)
	return nil
}

// Search 执行相似度搜索
func (v *VectorStore) Search(ctx context.Context, query string, topK int) ([]schema.Document, error) {
	fmt.Printf("🔍 [DEBUG] 检索查询: '%s', TopK: %d\n", query, topK)
	
	docs, err := v.store.SimilaritySearch(ctx, query, topK)
	if err != nil {
		return nil, fmt.Errorf("检索失败: %w", err)
	}
	
	fmt.Printf("📊 [DEBUG] 检索到 %d 个结果\n", len(docs))
	for i, doc := range docs {
		preview := doc.PageContent
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("   [%d] 相似度分数: %.4f, 内容预览: %s\n", i+1, doc.Score, preview)
	}
	
	return docs, nil
}
