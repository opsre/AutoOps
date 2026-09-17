// Markdown 解析器，将大模型返回的 Markdown 文本解析为结构化内容
package service

import (
	"dodevops-api/api/ai/model"
	"regexp"
	"strings"
)

// ParseMarkdown 解析 Markdown 文本为结构化内容块
func ParseMarkdown(text string) []model.ContentBlock {
	var blocks []model.ContentBlock

	// 按行分割
	lines := strings.Split(text, "\n")

	var i int
	for i < len(lines) {
		line := lines[i]

		// 1. 检测代码块 ```
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			codeBlock, endIdx := extractCodeBlock(lines, i)
			if codeBlock != nil {
				blocks = append(blocks, *codeBlock)
				i = endIdx + 1
				continue
			}
		}

		// 2. 检测标题 #
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			heading := extractHeading(line)
			if heading != nil {
				blocks = append(blocks, *heading)
				i++
				continue
			}
		}

		// 3. 检测列表项 - 或 * 或 数字.
		if isListItem(line) {
			listBlock, endIdx := extractList(lines, i)
			if listBlock != nil {
				blocks = append(blocks, *listBlock)
				i = endIdx
				continue
			}
		}

		// 4. 普通文本（合并连续的文本行）
		textBlock, endIdx := extractTextBlock(lines, i)
		if textBlock != nil && strings.TrimSpace(textBlock.Content) != "" {
			blocks = append(blocks, *textBlock)
			i = endIdx
			continue
		}

		i++
	}

	return blocks
}

// extractCodeBlock 提取代码块
func extractCodeBlock(lines []string, startIdx int) (*model.ContentBlock, int) {
	firstLine := strings.TrimSpace(lines[startIdx])

	// 提取语言标识（如 ```bash, ```go 等）
	language := strings.TrimPrefix(firstLine, "```")
	language = strings.TrimSpace(language)
	if language == "" {
		language = "plaintext"
	}

	// 查找代码块结束位置
	var codeLines []string
	endIdx := startIdx + 1

	for endIdx < len(lines) {
		if strings.HasPrefix(strings.TrimSpace(lines[endIdx]), "```") {
			break
		}
		codeLines = append(codeLines, lines[endIdx])
		endIdx++
	}

	return &model.ContentBlock{
		Type:     "code",
		Content:  strings.Join(codeLines, "\n"),
		Language: language,
	}, endIdx
}

// extractHeading 提取标题
func extractHeading(line string) *model.ContentBlock {
	trimmed := strings.TrimSpace(line)

	// 计算 # 的数量
	level := 0
	for _, ch := range trimmed {
		if ch == '#' {
			level++
		} else {
			break
		}
	}

	if level == 0 || level > 6 {
		return nil
	}

	// 提取标题内容
	content := strings.TrimSpace(trimmed[level:])

	return &model.ContentBlock{
		Type:    "heading",
		Content: content,
		Level:   level,
	}
}

// isListItem 检查是否为列表项
func isListItem(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 无序列表: - 或 *
	if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
		return true
	}

	// 有序列表: 数字.
	matched, _ := regexp.MatchString(`^\d+\.\s+`, trimmed)
	return matched
}

// extractList 提取列表（合并连续的列表项）
func extractList(lines []string, startIdx int) (*model.ContentBlock, int) {
	var listItems []string
	i := startIdx

	for i < len(lines) {
		line := lines[i]
		if !isListItem(line) && strings.TrimSpace(line) != "" {
			break
		}

		if strings.TrimSpace(line) != "" {
			listItems = append(listItems, strings.TrimSpace(line))
		}
		i++
	}

	if len(listItems) == 0 {
		return nil, startIdx
	}

	return &model.ContentBlock{
		Type:    "list",
		Content: strings.Join(listItems, "\n"),
	}, i
}

// extractTextBlock 提取文本块（合并连续的普通文本）
func extractTextBlock(lines []string, startIdx int) (*model.ContentBlock, int) {
	var textLines []string
	i := startIdx

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// 遇到特殊标记（代码块、标题、列表）则停止
		if strings.HasPrefix(trimmed, "```") ||
			strings.HasPrefix(trimmed, "#") ||
			isListItem(line) {
			break
		}

		// 空行也停止（分段）
		if trimmed == "" && len(textLines) > 0 {
			break
		}

		if trimmed != "" {
			textLines = append(textLines, trimmed)
		}
		i++
	}

	if len(textLines) == 0 {
		return nil, startIdx + 1
	}

	return &model.ContentBlock{
		Type:    "text",
		Content: strings.Join(textLines, " "),
	}, i
}
