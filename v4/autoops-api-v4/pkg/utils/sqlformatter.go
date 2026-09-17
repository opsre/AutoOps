package utils

import (
	"regexp"
	"strings"
)

// FormatSQL 格式化 SQL 语句，添加换行和缩进
func FormatSQL(sql string) string {
	if sql == "" {
		return ""
	}

	// 移除多余的空格，但保留单个空格
	sql = strings.TrimSpace(sql)
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")

	// 特殊处理 CASE 语句
	sql = formatCaseStatements(sql)

	// 处理主要关键字换行
	sql = formatMainClauses(sql)

	// 处理 JOIN 子句（在分列之前）
	sql = formatJoinClauses(sql)

	// 处理 SELECT 列
	sql = formatSelectColumns(sql)

	// 处理 WHERE/AND/OR 条件
	sql = formatWhereConditions(sql)

	// 清理多余空行
	sql = regexp.MustCompile(`\n{3,}`).ReplaceAllString(sql, "\n")
	sql = strings.TrimSpace(sql)

	return sql
}

// formatCaseStatements 格式化 CASE 语句
func formatCaseStatements(sql string) string {
	// 匹配 CASE ... END 块
	caseRegex := regexp.MustCompile(`(?i)\bCASE\b(.*?)\bEND\b`)

	return caseRegex.ReplaceAllStringFunc(sql, func(match string) string {
		// 格式化 CASE 内部
		result := match

		// WHEN 换行并缩进
		result = regexp.MustCompile(`(?i)\s+WHEN\s+`).ReplaceAllString(result, "\n    WHEN ")

		// THEN 保持在同一行
		result = regexp.MustCompile(`(?i)\s+THEN\s+`).ReplaceAllString(result, " THEN ")

		// ELSE 换行并缩进
		result = regexp.MustCompile(`(?i)\s+ELSE\s+`).ReplaceAllString(result, "\n    ELSE ")

		// END 换行并缩进
		result = regexp.MustCompile(`(?i)\s+END\b`).ReplaceAllString(result, "\n  END")

		// CASE 开头
		result = regexp.MustCompile(`(?i)^CASE\s+`).ReplaceAllString(result, "CASE ")

		return result
	})
}

// formatMainClauses 格式化主要子句
func formatMainClauses(sql string) string {
	// 主要子句换行
	mainClauses := []string{
		"FROM", "WHERE", "GROUP BY", "HAVING",
		"ORDER BY", "LIMIT", "OFFSET",
		"UNION", "UNION ALL",
	}

	for _, clause := range mainClauses {
		pattern := regexp.MustCompile(`(?i)\s+` + regexp.QuoteMeta(clause) + `\b`)
		sql = pattern.ReplaceAllString(sql, "\n"+strings.ToUpper(clause)+" ")
	}

	return sql
}

// formatJoinClauses 格式化 JOIN 子句
func formatJoinClauses(sql string) string {
	// 各种 JOIN 类型（按长度降序，避免 LEFT JOIN 被拆成 LEFT 和 JOIN）
	joinTypes := []string{
		"LEFT OUTER JOIN", "RIGHT OUTER JOIN", "FULL OUTER JOIN",
		"LEFT JOIN", "RIGHT JOIN", "INNER JOIN",
		"OUTER JOIN", "CROSS JOIN",
	}

	for _, joinType := range joinTypes {
		pattern := regexp.MustCompile(`(?i)\s+` + regexp.QuoteMeta(joinType) + `\s+`)
		sql = pattern.ReplaceAllString(sql, "\n  "+strings.ToUpper(joinType)+" ")
	}

	// ON 子句不单独换行，保持在 JOIN 同一行
	return sql
}

// formatSelectColumns 格式化 SELECT 列
func formatSelectColumns(sql string) string {
	lines := strings.Split(sql, "\n")
	result := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 如果是 SELECT 行
		if regexp.MustCompile(`(?i)^SELECT\s+`).MatchString(trimmed) {
			// 提取 SELECT 后的内容
			content := regexp.MustCompile(`(?i)^SELECT\s+`).ReplaceAllString(trimmed, "")

			// 分割列（不在 CASE...END 内的逗号）
			columns := splitByCommaOutsideCase(content)

			if len(columns) > 1 {
				result = append(result, "SELECT")
				for i, col := range columns {
					col = strings.TrimSpace(col)
					// 处理 CASE 语句的内部换行，增加缩进
					if strings.Contains(col, "\n") {
						col = strings.ReplaceAll(col, "\n", "\n  ")
					}
					if i == len(columns)-1 {
						// 最后一列
						result = append(result, "  "+col)
					} else {
						// 非最后一列，添加逗号
						result = append(result, "  "+col+",")
					}
				}
			} else {
				result = append(result, trimmed)
			}
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// formatWhereConditions 格式化 WHERE 条件
func formatWhereConditions(sql string) string {
	lines := strings.Split(sql, "\n")
	result := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 如果是 WHERE 行
		if regexp.MustCompile(`(?i)^WHERE\s+`).MatchString(trimmed) {
			// 提取 WHERE 后的内容
			content := regexp.MustCompile(`(?i)^WHERE\s+`).ReplaceAllString(trimmed, "")

			// 分割 AND/OR
			conditions := splitByAndOr(content)

			if len(conditions) > 1 {
				result = append(result, "WHERE")
				for _, cond := range conditions {
					cond = strings.TrimSpace(cond)
					result = append(result, "  "+cond)
				}
			} else {
				result = append(result, trimmed)
			}
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// splitByCommaOutsideCase 按逗号分割，但忽略 CASE...END 内的逗号
func splitByCommaOutsideCase(s string) []string {
	var result []string
	var current strings.Builder
	caseDepth := 0

	i := 0
	for i < len(s) {
		// 检查是否是 CASE 开始
		if i+4 <= len(s) && strings.ToUpper(s[i:i+4]) == "CASE" {
			if i+4 >= len(s) || !isLetterOrDigit(s[i+4]) {
				caseDepth++
				current.WriteString(s[i : i+4])
				i += 4
				continue
			}
		}

		// 检查是否是 END 结束
		if i+3 <= len(s) && strings.ToUpper(s[i:i+3]) == "END" {
			if i+3 >= len(s) || !isLetterOrDigit(s[i+3]) {
				if caseDepth > 0 {
					caseDepth--
				}
				current.WriteString(s[i : i+3])
				i += 3
				continue
			}
		}

		// 如果是逗号且不在 CASE 内
		if s[i] == ',' && caseDepth == 0 {
			result = append(result, current.String())
			current.Reset()
			i++
			continue
		}

		current.WriteByte(s[i])
		i++
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// splitByAndOr 按 AND/OR 分割
func splitByAndOr(s string) []string {
	var result []string
	var current strings.Builder

	i := 0
	for i < len(s) {
		// 检查是否是 AND
		if i+3 <= len(s) && strings.ToUpper(s[i:i+3]) == "AND" {
			if (i == 0 || !isLetterOrDigit(s[i-1])) && (i+3 >= len(s) || !isLetterOrDigit(s[i+3])) {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
				current.WriteString("AND ")
				i += 3
				// 跳过后续空格
				for i < len(s) && s[i] == ' ' {
					i++
				}
				continue
			}
		}

		// 检查是否是 OR
		if i+2 <= len(s) && strings.ToUpper(s[i:i+2]) == "OR" {
			if (i == 0 || !isLetterOrDigit(s[i-1])) && (i+2 >= len(s) || !isLetterOrDigit(s[i+2])) {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
				current.WriteString("OR ")
				i += 2
				// 跳过后续空格
				for i < len(s) && s[i] == ' ' {
					i++
				}
				continue
			}
		}

		current.WriteByte(s[i])
		i++
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// isLetterOrDigit 检查字符是否是字母或数字
func isLetterOrDigit(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
