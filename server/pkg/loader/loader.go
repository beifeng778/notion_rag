// pkg/loader/loader.go
package loader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

// DocumentLoader 文档加载器
type DocumentLoader struct {
	DocsDir   string
	ChunkSize int
	Overlap   int
}

// NewDocumentLoader 创建新的文档加载器
func NewDocumentLoader(docsDir string, chunkSize, overlap int) *DocumentLoader {
	return &DocumentLoader{
		DocsDir:   docsDir,
		ChunkSize: chunkSize,
		Overlap:   overlap,
	}
}

// LoadAndSplit 加载并分块所有文档
func (l *DocumentLoader) LoadAndSplit() ([]schema.Document, error) {
	files, err := filepath.Glob(filepath.Join(l.DocsDir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("读取文档目录失败: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("未找到任何 Markdown 文件")
	}

	var allChunks []schema.Document
	for _, file := range files {
		chunks, err := l.loadFile(file)
		if err != nil {
			log.Printf("警告：加载文件 %s 失败: %v", file, err)
			continue
		}
		fmt.Printf(" [DEBUG] 文件 %s: %d 个文本块\n", filepath.Base(file), len(chunks))
		allChunks = append(allChunks, chunks...)
	}

	return allChunks, nil
}

// loadFile 加载单个文件
func (l *DocumentLoader) loadFile(file string) ([]schema.Document, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// 分块处理
	return l.splitDocuments([]string{string(content)})
}

// splitDocuments 对文档进行分块
func (l *DocumentLoader) splitDocuments(rawContents []string) ([]schema.Document, error) {
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(l.ChunkSize),
		textsplitter.WithChunkOverlap(l.Overlap),
		textsplitter.WithSeparators([]string{"\n\n", "\n", " ", ""}),
	)

	var allChunks []schema.Document
	for _, content := range rawContents {
		// 先用 Markdown splitter 初步处理
		mdSplitter := textsplitter.NewMarkdownTextSplitter()
		mdDocs, err := mdSplitter.SplitText(content)
		if err != nil {
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
	return allChunks, nil
}
