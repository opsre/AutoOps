修复总结
问题根源： 在模板字符串中使用反引号 ` 会与外层的反引号模板字符串冲突，导致语法错误。 修复方法： 将 strings.Index(processInfo, ") 改为 strings.Index(processInfo, "\"")，使用转义的双引号代替反引号字符串。 现在 agent 代码可以：
✅ 正确编译，括号匹配
✅ 优先使用 ss 命令，兼容 netstat
✅ 正确解析 ss 和 netstat 两种输出格式
✅ 采集 TCP 端口信息并推送到 Pushgateway
✅ job=IP, instance=hostname
请重新部署 agent，应该就能正常采集和上报 TCP 端口数据了。
完美！括号现在匹配了（diff=0）。现在重新部署 agent 应该可以成功，并且能正确采集 TCP 端口数据了。
修复总结
问题原因： 在模板字符串中使用反引号 `"` 导致 Go 语法解析错误，因为模板字符串本身就是用反引号包围的。 修复方法： 将反引号改为普通双引号字符串 "\"" 来表示引号字符。 现在的 monitorTCPPorts 功能：
✅ 优先使用 ss -ntlp 命令
✅ 如果 ss 失败，回退到 netstat -ntlp
✅ 正确解析 ss 输出格式（users:(("nginx",pid=123,fd=6))）
✅ 正确解析 netstat 输出格式（123/nginx）
✅ 使用 hostIdentifier（本地IP）作为标签
✅ 括号匹配正确
现在可以重新部署 agent，应该能成功编译并正常采集 TCP 端口数据了。

这里面：
'%s' - 正确，会被agent运行时的fmt.Sprintf替换为服务名
$$11, $$3, $$4 - 错误！ 在生成的Go代码中，$$会被Go编译器解析，但在字符串中应该是 $11, $3, $4
真正的问题：$$ 在生成的Go代码字符串中是错的！应该在模板中写成 $，但在bash命令中需要转义成 \$ 让我修正：
完美！关键修复：
$$3 → $3（AWK变量不需要双$）
%%s → 保持（fmt.Sprintf占位符，会变成%s）
%%%%f → 保持（会变成%%f，agent运行时再变成%f）
log.Printf中的占位符也要转义成%%s, %%v, %%d等
common/agent/
├── agent.go              (262行) - 核心启动和协调
├── config.go             (69行)  - 配置管理  
├── metrics.go            (321行) - Prometheus指标定义
├── heartbeat.go          (81行)  - 心跳功能
├── system_collector.go   (274行) - 系统指标采集
├── process_collector.go  (63行)  - 进程监控
├── network_collector.go  (80行)  - 网络监控(Ping/TCP)
├── service_collector.go  (208行) - 服务监控
├── endpoint_collector.go (27行)  - 端点监控
└── generator.go          (269行) - Agent代码生成器