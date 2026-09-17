# 🔥 自动补全所有缺失的API描述
cd scripts && make fix-api-desc

这次改动会影响新增模型，自动创建表的功能吗
pkg/db/migrate.go
不会，功能保留了。

现在的规则是：

已有库、模型没变化：跳过全量 AutoMigrate
你新增模型、字段、gorm tag 变化：模型签名会变化，仍然自动执行 AutoMigrate
空库或新库：也会正常自动建表
也就是说，原来的“新增 model 后自动建表”还在，只是从“每次启动都全量扫描所有表”改成了“检测到模型结构变化才迁移”。

核心逻辑在 pkg/db/migrate.go。
你只要把新模型加进 models 列表，服务启动时还是会自动建表。

go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o devops main.go
lsof -i:8000 | grep LISTEN | awk '{print $2}' | xargs kill -9
 /Users/apple/go/bin/swag  init
 go run main.go

重定向 URL
http://192.168.3.7:8080/api/v1/tool/feishu/oauth/callback

# 激活一年
go run scripts/generate_license.go \
  --private-key ./scripts/license_private.pem \
  --customer "ops" \
  --machine-code "04b69de2bbf7387817e8c0b38ff4e9f56098324b399688083549523122cb5d9b" \
  --license-type annual \
  --issued-at "2026-03-23 00:00:00" \
  --expire-at "2027-03-23 23:59:59"

# 永久激活
go run scripts/generate_license.go \
  --private-key ./scripts/license_private.pem \
  --customer "Lockin" \
  --machine-code "激活码" \
  --license-type permanent
