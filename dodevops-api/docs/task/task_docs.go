package docstask

const TaskPaths = `
        "/api/v1/task/add": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "创建新的任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "创建任务",
                "parameters": [
                    {
                        "description": "任务信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CreateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "创建Ansible任务（1=手动，2=Git导入）。K8s部署任务请使用专门的K8s创建接口",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "创建Ansible任务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "任务类型(1=手动，2=Git导入)",
                        "name": "type",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "主机分组JSON",
                        "name": "hostGroups",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Git仓库地址(type=2时必填)",
                        "name": "gitRepo",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "全局变量JSON",
                        "name": "variables",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "playbook文件(type=1时上传，支持多文件)",
                        "name": "playbooks",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "roles目录(type=1时上传)",
                        "name": "roles",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "额外变量(JSON/YAML字符串)",
                        "name": "extra_vars",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "命令行参数",
                        "name": "cli_args",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "是否使用配置中心(0=否，1=是)",
                        "name": "use_config",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "Inventory配置ID",
                        "name": "inventory_config_id",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "全局变量配置ID",
                        "name": "global_vars_config_id",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "额外变量配置ID",
                        "name": "extra_vars_config_id",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "命令行参数配置ID",
                        "name": "cli_args_config_id",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Cron表达式(周期任务必填)",
                        "name": "cron_expr",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "是否为周期任务(0=否, 1=是)",
                        "name": "is_recurring",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Playbook文件路径列表(JSON数组字符串, type=2时可选)",
                        "name": "playbook_paths",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "视图ID",
                        "name": "view_id",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.TaskAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/query/name": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务名称进行模糊查询",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "根据名称模糊查询Ansible任务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务名称（支持模糊查询）",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.TaskAnsible"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/query/type": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务类型查询（1=手动，2=Git导入，3=K8s部署）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "根据类型查询Ansible任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务类型（1=手动，2=Git导入，3=K8s部署）",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.TaskAnsible"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/query": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "多条件查询Ansible任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "多条件查询Ansible任务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务名称（支持模糊查询）",
                        "name": "name",
                        "in": "query",
                        "required": false
                    },
                    {
                        "type": "integer",
                        "description": "任务类型（1=手动，2=Git导入，3=K8s部署）",
                        "name": "type",
                        "in": "query",
                        "required": false
                    },
                    {
                        "type": "string",
                        "description": "视图名称",
                        "name": "viewName",
                        "in": "query",
                        "required": false
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.TaskAnsible"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取Ansible任务详情",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取Ansible任务详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.TaskAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "修改Ansible任务基本信息和配置（运行中任务不可修改）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "修改Ansible任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "修改任务请求",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UpdateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.TaskAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除指定的Ansible任务（级联删除关联的子任务）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "删除Ansible任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/{id}/log/{work_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "通过SSE协议实时获取Ansible任务执行日志",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/event-stream"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取Ansible任务日志(SSE)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "子任务ID",
                        "name": "work_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SSE格式的实时日志",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/{id}/start": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "启动指定的Ansible任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "启动Ansible任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansiblelist": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取Ansible任务列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取Ansible任务列表",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/controller.ListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/{id}/history": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取任务的历史执行记录列表，支持分页",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取任务历史记录列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/history/{history_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取任务的历史执行详情，包含每个主机的执行日志",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取任务历史记录详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "历史记录ID",
                        "name": "history_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.TaskAnsibleHistory"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/{id}/history/{history_id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除指定的任务历史记录及关联的日志文件",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "删除任务历史记录",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "历史记录ID",
                        "name": "history_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            },
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取任务的历史执行详情，包含每个主机的执行日志",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取任务历史记录详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "历史记录ID",
                        "name": "history_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.TaskAnsibleHistory"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/history/detail/task/{task_id}/work/{work_id}/history/{history_id}/log": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID、WORKID和HistoryID获取历史任务日志",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取历史记录日志内容(通过详细信息)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "子任务ID",
                        "name": "work_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "历史记录ID",
                        "name": "history_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/ansible/history/work/{work_history_id}/log": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取指定子任务历史记录的日志内容",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取历史记录日志内容",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "子任务历史记录ID",
                        "name": "work_history_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "删除任务",
                "parameters": [
                    {
                        "description": "任务ID请求",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TaskIDRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/execution-info": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取任务的执行次数和下次执行时间",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "获取任务执行信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object",
                                            "additionalProperties": true
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/get": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID查询任务详情",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据ID查询任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/k8s": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "创建K8s集群部署任务",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "创建K8s部署任务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "任务描述",
                        "name": "description",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "K8s集群名称",
                        "name": "cluster_name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "K8s集群版本",
                        "name": "cluster_version",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "部署模式(1=单节点,2=集群)",
                        "name": "deployment_mode",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Master节点主机ID数组JSON",
                        "name": "master_host_ids",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Worker节点主机ID数组JSON",
                        "name": "worker_host_ids",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "ETCD节点主机ID数组JSON",
                        "name": "etcd_host_ids",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "启用组件数组JSON",
                        "name": "enabled_components",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "私有仓库地址",
                        "name": "private_registry",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "仓库用户名",
                        "name": "registry_username",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "仓库密码",
                        "name": "registry_password",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.TaskAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/task/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取任务列表，支持分页和条件查询",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "获取任务列表",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "任务状态(1=等待中,2=运行中,3=成功,4=异常)",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/list-with-details": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取任务列表，支持分页和条件查询，包含模板和主机的详细信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "获取任务列表（包含关联信息）",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "任务状态(1=等待中,2=运行中,3=成功,4=异常)",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/next-execution": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据cron表达式计算下次执行时间",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "计算下次执行时间",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cron表达式",
                        "name": "cron",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/query/name": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务名称模糊查询任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据名称查询任务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/query/status": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务状态查询任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据状态查询任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务状态(1=等待中,2=运行中,3=成功,4=异常)",
                        "name": "status",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/query/type": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务类型查询任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据类型查询任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务类型(1=普通任务,2=定时任务,3=ansible任务,4=工作作业)",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/templates": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID获取关联模板信息及状态",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "获取任务模板及状态",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/task/update": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "更新任务信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "更新任务",
                "parameters": [
                    {
                        "description": "需要更新的任务字段(必须包含ID)",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/taskjob/log": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID和模板ID获取日志",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取任务日志",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "taskId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "templateId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/taskjob/start": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID启动任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "启动任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "taskId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/taskjob/status": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID和模板ID获取任务状态",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "获取任务状态",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "taskId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "templateId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/taskjob/stop": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据任务ID和模板ID停止任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务作业"
                ],
                "summary": "停止单个任务",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "taskId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "templateId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/add": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "创建新的任务模板",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "创建任务模板",
                "parameters": [
                    {
                        "description": "任务模板信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TaskTemplate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/content/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取指定ID的任务模板脚本内容",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "获取脚本内容",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "脚本内容",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/template/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除指定ID的任务模板",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "删除模板",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/info/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取指定ID的任务模板",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据ID获取模板",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取所有任务模板列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "获取所有模板",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/query/name": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取名称包含指定字符串的任务模板列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据名称模糊查询模板",
                "parameters": [
                    {
                        "type": "string",
                        "description": "模板名称",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/query/type": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取指定类型的任务模板列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "根据类型查询模板",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "模板类型(1=shell, 2=python, 3=ansible)",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/template/update": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "更新任务模板",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务中心"
                ],
                "summary": "更新模板",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "模板ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "需要更新的模板字段",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TaskTemplate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/ws/task/ansible/{id}/log/{work_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "建立WebSocket连接实时推送任务执行日志",
                "tags": [
                    "任务作业"
                ],
                "summary": "通过WebSocket实时获取Ansible任务日志",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "子任务ID",
                        "name": "work_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/config/ansible": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "分页获取配置列表，支持按名称和类型过滤",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config配置中心"
                ],
                "summary": "获取Ansible配置列表",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "配置名称（模糊查询）",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "配置类型(1-inventory 2-global_vars 3-extra_vars 4-cli_args)",
                        "name": "type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dao.ListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "创建Inventory/Vars/Args等配置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config配置中心"
                ],
                "summary": "创建Ansible配置",
                "parameters": [
                    {
                        "description": "创建配置请求",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.CreateConfigRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.ConfigAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/config/ansible/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "根据ID获取配置详情",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config配置中心"
                ],
                "summary": "获取Ansible配置详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "配置ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.ConfigAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "更新配置内容",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config配置中心"
                ],
                "summary": "更新Ansible配置",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "配置ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "更新配置请求",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UpdateConfigRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/result.Result"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.ConfigAnsible"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除指定的配置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config配置中心"
                ],
                "summary": "删除Ansible配置",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "配置ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Result"
                        }
                    }
                }
            }
        }`

const TaskDefinitions = `

        "controller.CreateTaskRequest": {
            "type": "object",
            "required": [
                "host_ids",
                "name",
                "shell",
                "type"
            ],
            "properties": {
                "cron_expr": {
                    "type": "string"
                },
                "host_ids": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "shell": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "model.JenkinsTask": {
            "type": "object",
            "properties": {
                "color": {
                    "description": "任务颜色",
                    "type": "string"
                },
                "name": {
                    "description": "任务名称",
                    "type": "string"
                },
                "url": {
                    "description": "任务URL",
                    "type": "string"
                }
            }
        },
        "model.Task": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "cron_expr": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "end_time": {
                    "type": "string"
                },
                "execute_count": {
                    "type": "integer"
                },
                "host_ids": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "next_run_time": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "shell": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "task_count": {
                    "type": "integer"
                },
                "tasklog": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "model.TaskAnsible": {
            "type": "object",
            "properties": {
                "allHostIDs": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "errorMsg": {
                    "type": "string"
                },
                "gitRepo": {
                    "type": "string"
                },
                "globalVars": {
                    "type": "string"
                },
                "hostGroups": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "taskCount": {
                    "type": "integer"
                },
                "totalDuration": {
                    "type": "integer"
                },
                "type": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                },
                "works": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.TaskAnsibleWork"
                    }
                }
            }
        },
        "model.TaskAnsibleWork": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                },
                "endTime": {
                    "type": "string"
                },
                "entryFileName": {
                    "type": "string"
                },
                "entryFilePath": {
                    "type": "string"
                },
                "errorMsg": {
                    "type": "string"
                },
                "exitCode": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "logPath": {
                    "type": "string"
                },
                "startTime": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "task": {
                    "$ref": "#/definitions/model.TaskAnsible"
                },
                "taskID": {
                    "type": "integer"
                }
            }
        },
        "model.TaskAnsibleHistory": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "errorMsg": {
                    "type": "string"
                },
                "finishedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "operatorID": {
                    "type": "integer"
                },
                "operatorName": {
                    "type": "string"
                },
                "startedAt": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "taskAnsible": {
                    "$ref": "#/definitions/model.TaskAnsible"
                },
                "taskID": {
                    "type": "integer"
                },
                "totalDuration": {
                    "type": "integer"
                },
                "trigger": {
                    "type": "integer"
                },
                "uniqId": {
                    "type": "string"
                },
                "workHistories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.TaskAnsibleworkHistory"
                    }
                }
            }
        },
        "model.TaskAnsibleworkHistory": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "historyID": {
                    "type": "integer"
                },
                "hostName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "logPath": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "taskID": {
                    "type": "integer"
                },
                "workID": {
                    "type": "integer"
                }
            }
        },
        "model.TaskIDRequest": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.TaskStats": {
            "type": "object",
            "properties": {
                "failed": {
                    "description": "失败次数",
                    "type": "integer"
                },
                "success": {
                    "description": "成功次数",
                    "type": "integer"
                },
                "successRate": {
                    "description": "成功率",
                    "type": "number"
                },
                "total": {
                    "description": "任务执行总次数",
                    "type": "integer"
                }
            }
        },
        "model.TaskStatusResponse": {
            "type": "object",
            "properties": {
                "app_code": {
                    "description": "应用编码",
                    "type": "string"
                },
                "app_name": {
                    "description": "应用名称",
                    "type": "string"
                },
                "build_number": {
                    "description": "构建编号",
                    "type": "integer"
                },
                "duration": {
                    "description": "耗时(秒)",
                    "type": "integer"
                },
                "end_time": {
                    "description": "结束时间",
                    "type": "string"
                },
                "environment": {
                    "description": "环境名称",
                    "type": "string"
                },
                "error_message": {
                    "description": "错误信息",
                    "type": "string"
                },
                "log_url": {
                    "description": "日志URL",
                    "type": "string"
                },
                "progress": {
                    "description": "进度百分比(0-100)",
                    "type": "integer"
                },
                "start_time": {
                    "description": "开始时间",
                    "type": "string"
                },
                "status": {
                    "description": "任务状态: 1=未部署 2=部署中 3=成功 4=异常",
                    "type": "integer"
                },
                "status_text": {
                    "description": "状态文本",
                    "type": "string"
                },
                "task_id": {
                    "description": "任务ID",
                    "type": "integer"
                }
            }
        },
        "model.TaskTemplate": {
            "type": "object",
            "properties": {
                "content": {
                    "description": "任务内容",
                    "type": "string"
                },
                "createdAt": {
                    "description": "创建时间",
                    "type": "string"
                },
                "createdBy": {
                    "description": "创建人",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "description": "任务名称",
                    "type": "string"
                },
                "remark": {
                    "description": "备注信息",
                    "type": "string"
                },
                "type": {
                    "description": "1=shell模板, 2=python模板, 3=ansible模板",
                    "type": "integer"
                },
                "updatedAt": {
                    "description": "更新时间",
                    "type": "string"
                },
                "updatedBy": {
                    "description": "更新人",
                    "type": "string"
                }
            }
        },
        "model.ConfigAnsible": {
            "type": "object",
            "properties": {
                "content": {
                    "description": "内容：inventory为文本，vars/args为JSON",
                    "type": "string"
                },
                "createdAt": {
                    "description": "创建时间",
                    "type": "string"
                },
                "createdBy": {
                    "description": "创建人",
                    "type": "string"
                },
                "id": {
                    "description": "主键ID",
                    "type": "integer"
                },
                "name": {
                    "description": "配置名称",
                    "type": "string"
                },
                "remark": {
                    "description": "备注",
                    "type": "string"
                },
                "type": {
                    "description": "1-inventory 2-global_vars 3-extra_vars 4-cli_args",
                    "type": "integer"
                },
                "updatedAt": {
                    "description": "更新时间",
                    "type": "string"
                },
                "updatedBy": {
                    "description": "更新人",
                    "type": "string"
                }
            }
        },
		"service.CreateConfigRequest": {
            "type": "object",
            "required": [
                "content",
                "name",
                "type"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "type": {
                    "description": "1-inventory 2-global_vars 3-extra_vars 4-cli_args",
                    "type": "integer"
                }
            }
        },
        "service.UpdateConfigRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "remark": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "service.UpdateTaskRequest": {
            "type": "object",
            "properties": {
                "cliArgs": {
                    "type": "string",
					"description": "命令行参数"
                },
                "cliArgsConfigId": {
                    "type": "integer",
					"description": "命令行参数配置ID"
                },
                "extraVars": {
                    "type": "string",
					"description": "额外变量"
                },
                "extraVarsConfigId": {
                    "type": "integer",
					"description": "额外变量配置ID"
                },
                "gitRepo": {
                    "type": "string",
					"description": "Git代码库地址"
                },
                "globalVarsConfigId": {
                    "type": "integer",
					"description": "全局变量配置ID"
                },
                "hostGroups": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        }
                    },
					"description": "主机组"
                },
                "inventoryConfigId": {
                    "type": "integer",
					"description": "Inventory配置ID"
                },
                "name": {
                    "type": "string",
					"description": "任务名称"
                },
                "playbookPaths": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
					"description": "剧本路径列表"
                },
                "useConfig": {
                    "type": "integer",
					"description": "是否使用配置:0-否,1-是"
                },
                "variables": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
					"description": "任务全局变量"
                },
				"cronExpr": {
					"type": "string",
					"description": "定时表达式"
				},
				"isRecurring": {
					"type": "integer",
					"description": "是否周期性任务:0-否,1-是"
				},
				"viewId": {
					"type": "integer",
					"description": "视图ID"
				}
            }
        }`
