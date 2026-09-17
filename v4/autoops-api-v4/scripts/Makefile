# API描述检查相关命令

.PHONY: check-api-desc
check-api-desc:
	@echo "🔍 检查API描述覆盖率..."
	@cd .. && go run scripts/check_api_descriptions.go

.PHONY: check-api-missing
check-api-missing:
	@echo "⚠️  查看缺失的API描述..."
	@cd .. && go run scripts/check_api_descriptions.go 2>/dev/null | grep -A 500 "⚠️  发现"

.PHONY: check-api-stats
check-api-stats:
	@echo "📊 查看统计信息..."
	@cd .. && go run scripts/check_api_descriptions.go 2>/dev/null | grep -E "(扫描到|共配置|覆盖率)"

.PHONY: fix-api-desc
fix-api-desc:
	@echo "🔧 自动补全API描述..."
	@cd .. && go run scripts/auto_fix_api_descriptions.go

.PHONY: help
help:
	@echo "可用命令:"
	@echo "  make check-api-desc      - 完整检查API描述覆盖率"
	@echo "  make check-api-missing   - 仅显示缺失的API"
	@echo "  make check-api-stats     - 仅显示统计信息"
	@echo "  make fix-api-desc        - 🔥 自动补全所有缺失的API描述"
