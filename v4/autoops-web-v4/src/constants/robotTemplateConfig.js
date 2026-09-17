/**
 * 机器人告警模板配置
 * 为不同平台提供预设模板
 */

export const ROBOT_TEMPLATE_CONFIG = {
  // 钉钉机器人模板
  dingtalk: {
    name: '钉钉机器人',
    type: 'dingtalk',
    description: '钉钉机器人支持Markdown格式，消息最大4000字节',
    templates: [
      {
        name: 'Markdown模板（推荐）',
        template: JSON.stringify({
          msgtype: 'markdown',
          markdown: {
            title: '{{if .IsRecovered}}✅ 恢复通知{{else}}⚠️ 告警通知{{end}}',
            text: '### {{if .IsRecovered}}✅ 恢复通知{{else}}⚠️ 告警通知{{end}}\n\n---\n- **报警类型**: ${rule_name}\n- **报警指纹**: ${fingerprint}\n- **报警等级**: ${severity}\n- **报警主机**: ${labels.node_name}\n- **容器名称**: ${labels.pod}\n- **业务环境**: ${labels.namespace}\n- **触发时间**: ${first_trigger_time_format}{{if .IsRecovered}}\n- **恢复时间**: ${recover_time_format}{{end}}\n- **值班人员**: ${duty_user}\n\n---\n\n**报警事件**: ${annotations}\n\n[查看详情](http://localhost:3000/faultCenter/detail/${faultCenterId})'
          }
        }),
        preview: '### ⚠️ 告警通知\n---\n- 报警类型: CPU使用率过高\n- 报警等级: critical\n...'
      },
      {
        name: '简洁文本模板',
        template: JSON.stringify({
          msgtype: 'text',
          text: {
            content: '【{{if .IsRecovered}}恢复通知{{else}}告警通知{{end}}】\n规则：${rule_name}\n等级：${severity}\n主机：${labels.node_name}\n容器：${labels.pod}\n环境：${labels.namespace}\n时间：${first_trigger_time_format}\n值班：${duty_user}\n详情：http://localhost:3000/faultCenter/detail/${faultCenterId}'
          }
        }),
        preview: '【告警通知】\n规则：CPU使用率过高\n等级：critical\n...'
      }
    ],
    variables: [
      { var: '${rule_name}', desc: '报警规则名称', example: 'CPU使用率过高' },
      { var: '${fingerprint}', desc: '报警指纹', example: 'abc123def456' },
      { var: '${severity}', desc: '报警等级', example: 'critical/warning/info' },
      { var: '${labels.node_name}', desc: '报警主机', example: 'node-01' },
      { var: '${labels.pod}', desc: '容器名称', example: 'web-server-pod' },
      { var: '${labels.namespace}', desc: '业务环境', example: 'production' },
      { var: '${first_trigger_time_format}', desc: '触发时间', example: '2024-01-29 14:30:00' },
      { var: '${recover_time_format}', desc: '恢复时间', example: '2024-01-29 15:00:00' },
      { var: '${duty_user}', desc: '值班人员', example: '张三' },
      { var: '${annotations}', desc: '报警详情', example: 'CPU使用率超过90%' },
      { var: '${faultCenterId}', desc: '故障ID', example: '12345' }
    ],
    conditionals: [
      { syntax: '{{if .IsRecovered}}', desc: '判断是否恢复', usage: '{{if .IsRecovered}}已恢复{{else}}告警中{{end}}' }
    ],
    notes: [
      '钉钉机器人消息最大长度为4000字节',
      'Markdown格式支持有限，建议使用简单格式',
      '如需@某人，需在群内添加机器人并配置关键词'
    ]
  },

  // 飞书机器人模板
  feishu: {
    name: '飞书机器人',
    type: 'feishu',
    description: '飞书机器人支持多种消息格式，推荐使用交互式卡片',
    templates: [
      {
        name: '交互式卡片（推荐）',
        template: JSON.stringify({
          msg_type: 'interactive',
          card: {
            header: {
              title: {
                tag: 'plain_text',
                content: '【{{if .IsRecovered}}✅ 恢复通知{{else}}🚨 告警通知{{end}}】${alert_title}'
              },
              template: '{{if .IsRecovered}}green{{else}}red{{end}}'
            },
            elements: [
              {
                tag: 'div',
                fields: [
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**规则名称**\n${rule_name}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**报警指纹**\n${fingerprint}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**报警等级**\n${severity}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**报警主机**\n${target_name}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**容器名称**\n${labels.pod}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**业务环境**\n${labels.namespace}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**开始时间**\n${first_trigger_time_format}'
                    }
                  },
                  {
                    is_short: true,
                    text: {
                      tag: 'lark_md',
                      content: '**值班人员**\n${duty_user}'
                    }
                  }
                ]
              },
              {
                tag: 'action',
                actions: [
                  {
                    tag: 'button',
                    text: {
                      tag: 'plain_text',
                      content: '查看详情'
                    },
                    type: 'primary',
                    url: 'http://localhost:3000/faultCenter/detail/${faultCenterId}'
                  }
                ]
              }
            ]
          }
        }),
        preview: '【告警通知】ops-内存使用率 > 80%'
      },
      {
        name: '文本模板',
        template: JSON.stringify({
          msg_type: 'text',
          content: '【{{if .IsRecovered}}恢复通知{{else}}告警通知{{end}}】${alert_title}\n\n规则名称：${rule_name}\n报警指纹：${fingerprint}\n报警等级：${severity}\n报警主机：${target_name}\n容器名称：${labels.pod}\n业务环境：${labels.namespace}\n开始时间：${first_trigger_time_format}\n值班人员：${duty_user}\n详情：http://localhost:3000/faultCenter/detail/${faultCenterId}'
        }),
        preview: '【告警通知】ops-内存使用率 > 80%\n\n规则名称：...'
      }
    ],
    variables: [
      { var: '${rule_name}', desc: '报警规则名称', example: 'CPU使用率过高' },
      { var: '${alert_title}', desc: '告警标题', example: 'ops-内存使用率 > 80%' },
      { var: '${fingerprint}', desc: '报警指纹（唯一标识）', example: 'abc123def456' },
      { var: '${severity}', desc: '报警等级', example: 'critical/warning/info' },
      { var: '${target_name}', desc: '报警主机（含IP）', example: 'go-ops (172.20.236.121)' },
      { var: '${labels.pod}', desc: '容器/Pod名称', example: 'web-server-pod' },
      { var: '${labels.namespace}', desc: '业务环境/命名空间', example: 'production' },
      { var: '${first_trigger_time_format}', desc: '首次触发时间', example: '2024-01-29 14:30:00' },
      { var: '${recover_time_format}', desc: '恢复时间', example: '2024-01-29 15:00:00' },
      { var: '${duty_user}', desc: '值班人员', example: '张三' },
      { var: '${faultCenterId}', desc: '故障中心ID', example: '12345' }
    ],
    conditionals: [
      { syntax: '{{if .IsRecovered}}', desc: '判断是否已恢复', usage: '{{if .IsRecovered}}已恢复{{else}}告警中{{end}}' }
    ],
    notes: [
      '交互式卡片支持颜色区分：green（恢复/绿色）、red（告警/红色）',
      '支持按钮交互，用户可直接点击"查看详情"跳转到故障中心',
      '字段采用两列布局，信息清晰，易于阅读'
    ]
  },

  // 企业微信机器人模板
  wechat: {
    name: '企业微信机器人',
    type: 'wechat',
    description: '企业微信机器人支持Markdown格式',
    templates: [
      {
        name: 'Markdown模板（推荐）',
        template: JSON.stringify({
          msgtype: 'markdown',
          markdown: {
            content: '{{if .IsRecovered}}### <font color="info">✅ 恢复通知</font>{{else}}### <font color="warning">🚨 告警通知</font>{{end}}\n> **报警类型**: ${rule_name}\n> **报警指纹**: ${fingerprint}\n> **报警等级**: <font color="warning">${severity}</font>\n> **报警主机**: ${labels.node_name}\n> **容器名称**: ${labels.pod}\n> **业务环境**: ${labels.namespace}\n> **触发时间**: ${first_trigger_time_format}{{if .IsRecovered}}\n> **恢复时间**: ${recover_time_format}{{end}}\n> **值班人员**: ${duty_user}\n\n**报警事件**: ${annotations}\n\n[查看详情](http://localhost:3000/faultCenter/detail/${faultCenterId})'
          }
        }),
        preview: '### 🚨 告警通知\n> 报警类型: CPU使用率过高\n> 报警等级: critical\n...'
      },
      {
        name: '文本模板',
        template: JSON.stringify({
          msgtype: 'text',
          text: {
            content: '【{{if .IsRecovered}}恢复通知{{else}}告警通知{{end}}】\n\n报警类型：${rule_name}\n报警等级：${severity}\n报警主机：${labels.node_name}\n容器名称：${labels.pod}\n业务环境：${labels.namespace}\n触发时间：${first_trigger_time_format}\n值班人员：${duty_user}\n\n报警事件：${annotations}\n\n详情：http://localhost:3000/faultCenter/detail/${faultCenterId}'
          }
        }),
        preview: '【告警通知】\n\n报警类型：CPU使用率过高\n...'
      }
    ],
    variables: [
      { var: '${rule_name}', desc: '报警规则名称', example: 'CPU使用率过高' },
      { var: '${fingerprint}', desc: '报警指纹', example: 'abc123def456' },
      { var: '${severity}', desc: '报警等级', example: 'critical/warning/info' },
      { var: '${labels.node_name}', desc: '报警主机', example: 'node-01' },
      { var: '${labels.pod}', desc: '容器名称', example: 'web-server-pod' },
      { var: '${labels.namespace}', desc: '业务环境', example: 'production' },
      { var: '${first_trigger_time_format}', desc: '触发时间', example: '2024-01-29 14:30:00' },
      { var: '${recover_time_format}', desc: '恢复时间', example: '2024-01-29 15:00:00' },
      { var: '${duty_user}', desc: '值班人员', example: '张三' },
      { var: '${annotations}', desc: '报警详情', example: 'CPU使用率超过90%' },
      { var: '${faultCenterId}', desc: '故障ID', example: '12345' }
    ],
    conditionals: [
      { syntax: '{{if .IsRecovered}}', desc: '判断是否恢复', usage: '{{if .IsRecovered}}已恢复{{else}}告警中{{end}}' }
    ],
    notes: [
      '企业微信支持font标签着色：info(绿), warning(橙), comment(灰)',
      '支持@所有人：在content中添加 <@all>'
    ]
  },

  // 邮件通知模板
  email: {
    name: '邮件通知',
    type: 'email',
    description: '邮件支持HTML格式',
    templates: [
      {
        name: 'HTML邮件模板（推荐）',
        template: JSON.stringify({
          subject: '{{if .IsRecovered}}[已恢复]{{else}}[告警]{{end}} ${rule_name}',
          body: '{{if .IsRecovered}}<div style="padding:20px;background:#f0f9ff;border-left:4px solid #0ea5e9;"><h2 style="color:#0ea5e9;margin:0 0 10px 0;">✅ 恢复通知</h2></div>{{else}}<div style="padding:20px;background:#fef2f2;border-left:4px solid #ef4444;"><h2 style="color:#ef4444;margin:0 0 10px 0;">🚨 告警通知</h2></div>{{end}}<div style="padding:20px;background:#f9fafb;"><table style="width:100%;border-collapse:collapse;"><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>🤖 报警类型</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${rule_name}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>🫧 报警指纹</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${fingerprint}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>📌 报警等级</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${severity}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>🖥 报警主机</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${labels.node_name}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>🧚 容器名称</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${labels.pod}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>☘️ 业务环境</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${labels.namespace}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>🕘 触发时间</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${first_trigger_time_format}</td></tr>{{if .IsRecovered}}<tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>🕐 恢复时间</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${recover_time_format}</td></tr>{{end}}<tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>👤 值班人员</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${duty_user}</td></tr><tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;"><strong>📝 报警事件</strong></td><td style="padding:8px;border-bottom:1px solid #e5e7eb;">${annotations}</td></tr></table></div><div style="padding:20px;text-align:center;"><a href="http://localhost:3000/faultCenter/detail/${faultCenterId}" style="display:inline-block;padding:10px 30px;background:#3b82f6;color:#fff;text-decoration:none;border-radius:5px;">查看详情</a></div>'
        }),
        preview: 'HTML格式邮件，包含样式和表格'
      },
      {
        name: '纯文本邮件模板',
        template: JSON.stringify({
          subject: '{{if .IsRecovered}}[已恢复]{{else}}[告警]{{end}} ${rule_name}',
          body: '{{if .IsRecovered}}========== ✅ 恢复通知 ==========\n{{else}}========== 🚨 告警通知 ==========\n{{end}}\n报警类型: ${rule_name}\n报警指纹: ${fingerprint}\n报警等级: ${severity}\n报警主机: ${labels.node_name}\n容器名称: ${labels.pod}\n业务环境: ${labels.namespace}\n触发时间: ${first_trigger_time_format}{{if .IsRecovered}}\n恢复时间: ${recover_time_format}{{end}}\n值班人员: ${duty_user}\n\n报警事件:\n${annotations}\n\n查看详情: http://localhost:3000/faultCenter/detail/${faultCenterId}'
        }),
        preview: '========== 🚨 告警通知 ==========\n报警类型: CPU使用率过高\n...'
      }
    ],
    variables: [
      { var: '${rule_name}', desc: '报警规则名称', example: 'CPU使用率过高' },
      { var: '${fingerprint}', desc: '报警指纹', example: 'abc123def456' },
      { var: '${severity}', desc: '报警等级', example: 'critical/warning/info' },
      { var: '${labels.node_name}', desc: '报警主机', example: 'node-01' },
      { var: '${labels.pod}', desc: '容器名称', example: 'web-server-pod' },
      { var: '${labels.namespace}', desc: '业务环境', example: 'production' },
      { var: '${first_trigger_time_format}', desc: '触发时间', example: '2024-01-29 14:30:00' },
      { var: '${recover_time_format}', desc: '恢复时间', example: '2024-01-29 15:00:00' },
      { var: '${duty_user}', desc: '值班人员', example: '张三' },
      { var: '${annotations}', desc: '报警详情', example: 'CPU使用率超过90%' },
      { var: '${faultCenterId}', desc: '故障ID', example: '12345' }
    ],
    conditionals: [
      { syntax: '{{if .IsRecovered}}', desc: '判断是否恢复', usage: '{{if .IsRecovered}}已恢复{{else}}告警中{{end}}' }
    ],
    notes: [
      'subject和body需同时配置',
      'HTML格式支持CSS样式',
      '建议使用内联样式以确保邮件客户端兼容性'
    ]
  },

  // 通用Webhook模板
  webhook: {
    name: '通用Webhook',
    type: 'webhook',
    description: '通用Webhook支持自定义JSON格式',
    templates: [
      {
        name: '标准JSON模板',
        template: JSON.stringify({
          alert_type: '${rule_name}',
          fingerprint: '${fingerprint}',
          severity: '${severity}',
          status: '{{if .IsRecovered}}resolved{{else}}firing{{end}}',
          labels: {
            node_name: '${labels.node_name}',
            pod: '${labels.pod}',
            namespace: '${labels.namespace}'
          },
          annotations: '${annotations}',
          trigger_time: '${first_trigger_time_format}',
          recover_time: '${recover_time_format}',
          duty_user: '${duty_user}',
          detail_url: 'http://localhost:3000/faultCenter/detail/${faultCenterId}'
        }),
        preview: '{\n  "alert_type": "CPU使用率过高",\n  "severity": "critical",\n  "status": "firing",\n  ...'
      },
      {
        name: '简化JSON模板',
        template: JSON.stringify({
          title: '${rule_name}',
          level: '${severity}',
          message: '${annotations}',
          host: '${labels.node_name}',
          time: '${first_trigger_time_format}',
          url: 'http://localhost:3000/faultCenter/detail/${faultCenterId}'
        }),
        preview: '{\n  "title": "CPU使用率过高",\n  "level": "critical",\n  ...'
      }
    ],
    variables: [
      { var: '${rule_name}', desc: '报警规则名称', example: 'CPU使用率过高' },
      { var: '${fingerprint}', desc: '报警指纹', example: 'abc123def456' },
      { var: '${severity}', desc: '报警等级', example: 'critical/warning/info' },
      { var: '${labels.node_name}', desc: '报警主机', example: 'node-01' },
      { var: '${labels.pod}', desc: '容器名称', example: 'web-server-pod' },
      { var: '${labels.namespace}', desc: '业务环境', example: 'production' },
      { var: '${first_trigger_time_format}', desc: '触发时间', example: '2024-01-29 14:30:00' },
      { var: '${recover_time_format}', desc: '恢复时间', example: '2024-01-29 15:00:00' },
      { var: '${duty_user}', desc: '值班人员', example: '张三' },
      { var: '${annotations}', desc: '报警详情', example: 'CPU使用率超过90%' },
      { var: '${faultCenterId}', desc: '故障ID', example: '12345' }
    ],
    conditionals: [
      { syntax: '{{if .IsRecovered}}', desc: '判断是否恢复', usage: '{{if .IsRecovered}}resolved{{else}}firing{{end}}' }
    ],
    notes: [
      '支持任意JSON格式',
      '可根据目标系统API要求自定义字段',
      'JSON必须合法，否则发送失败'
    ]
  },

  // Microsoft Teams模板
  teams: {
    name: 'Microsoft Teams',
    type: 'teams',
    description: 'Microsoft Teams使用MessageCard格式',
    templates: [
      {
        name: 'MessageCard模板（推荐）',
        template: JSON.stringify({
          '@type': 'MessageCard',
          '@context': 'https://schema.org/extensions',
          summary: '${rule_name}',
          themeColor: '{{if .IsRecovered}}2DC72D{{else}}E81123{{end}}',
          title: '{{if .IsRecovered}}✅ 恢复通知{{else}}🚨 告警通知{{end}}',
          sections: [
            {
              activityTitle: '${rule_name}',
              facts: [
                { name: '报警指纹', value: '${fingerprint}' },
                { name: '报警等级', value: '${severity}' },
                { name: '报警主机', value: '${labels.node_name}' },
                { name: '容器名称', value: '${labels.pod}' },
                { name: '业务环境', value: '${labels.namespace}' },
                { name: '触发时间', value: '${first_trigger_time_format}' }
              ],
              text: '${annotations}'
            }
          ],
          potentialAction: [
            {
              '@type': 'OpenUri',
              name: '查看详情',
              targets: [
                {
                  os: 'default',
                  uri: 'http://localhost:3000/faultCenter/detail/${faultCenterId}'
                }
              ]
            }
          ]
        }),
        preview: 'Teams卡片格式，包含按钮和格式化信息'
      }
    ],
    variables: [
      { var: '${rule_name}', desc: '报警规则名称', example: 'CPU使用率过高' },
      { var: '${fingerprint}', desc: '报警指纹', example: 'abc123def456' },
      { var: '${severity}', desc: '报警等级', example: 'critical/warning/info' },
      { var: '${labels.node_name}', desc: '报警主机', example: 'node-01' },
      { var: '${labels.pod}', desc: '容器名称', example: 'web-server-pod' },
      { var: '${labels.namespace}', desc: '业务环境', example: 'production' },
      { var: '${first_trigger_time_format}', desc: '触发时间', example: '2024-01-29 14:30:00' },
      { var: '${recover_time_format}', desc: '恢复时间', example: '2024-01-29 15:00:00' },
      { var: '${duty_user}', desc: '值班人员', example: '张三' },
      { var: '${annotations}', desc: '报警详情', example: 'CPU使用率超过90%' },
      { var: '${faultCenterId}', desc: '故障ID', example: '12345' }
    ],
    conditionals: [
      { syntax: '{{if .IsRecovered}}', desc: '判断是否恢复', usage: '{{if .IsRecovered}}2DC72D{{else}}E81123{{end}}' }
    ],
    notes: [
      '必须使用MessageCard格式',
      'themeColor使用十六进制颜色代码（不带#）',
      '支持potentialAction添加按钮'
    ]
  }
}

/**
 * 获取机器人类型的模板配置
 */
export const getRobotTemplateConfig = (type) => {
  return ROBOT_TEMPLATE_CONFIG[type] || null
}

/**
 * 获取所有支持的机器人类型
 */
export const getSupportedRobotTypes = () => {
  return Object.keys(ROBOT_TEMPLATE_CONFIG).map(key => ({
    value: key,
    label: ROBOT_TEMPLATE_CONFIG[key].name
  }))
}

/**
 * 获取指定类型的第一个模板
 */
export const getDefaultTemplate = (type) => {
  const config = getRobotTemplateConfig(type)
  return config?.templates?.[0]?.template || ''
}
