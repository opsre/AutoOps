#!/usr/bin/env python3
"""
API 描述自动更新工具 - 一键完成所有操作
功能：
1. 扫描所有 controller 文件，提取 Swagger 注解
2. 自动更新 middleware/apiDescriptions.go 文件
3. 生成统计报告
"""

import os
import re
import json
from pathlib import Path
from typing import List, Dict, Optional


class APIInfo:
    """API 信息类"""
    def __init__(self):
        self.tags: Optional[str] = None
        self.summary: Optional[str] = None
        self.description: Optional[str] = None
        self.method: Optional[str] = None
        self.router: Optional[str] = None
        self.file_path: Optional[str] = None
        self.line_number: Optional[int] = None

    def to_dict(self) -> Dict:
        return {
            'tags': self.tags,
            'summary': self.summary,
            'description': self.description,
            'method': self.method,
            'router': self.router,
            'file_path': self.file_path,
            'line_number': self.line_number
        }

    def is_valid(self) -> bool:
        return self.router is not None and self.method is not None


def extract_api_info_from_file(file_path: str) -> List[APIInfo]:
    """从 Go 文件中提取 API 信息"""
    apis = []
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except Exception as e:
        print(f"  ⚠️  读取失败: {file_path}")
        return apis

    current_api = None
    for line_num, line in enumerate(lines, start=1):
        line = line.strip()

        if line.startswith('// @Tags'):
            if current_api and current_api.is_valid():
                apis.append(current_api)
            current_api = APIInfo()
            current_api.file_path = file_path
            current_api.line_number = line_num
            current_api.tags = line.replace('// @Tags', '').strip()

        elif line.startswith('// @Summary') and current_api:
            current_api.summary = line.replace('// @Summary', '').strip()

        elif line.startswith('// @Description') and current_api:
            current_api.description = line.replace('// @Description', '').strip()

        elif line.startswith('// @Router') and current_api:
            match = re.search(r'@Router\s+(\S+)\s+\[(\w+)\]', line)
            if match:
                current_api.router = match.group(1)
                current_api.method = match.group(2).upper()

        elif line.startswith('func ') and current_api:
            if current_api.is_valid():
                apis.append(current_api)
            current_api = None

    if current_api and current_api.is_valid():
        apis.append(current_api)

    return apis


def scan_controllers(base_path: str = '.') -> List[APIInfo]:
    """扫描所有 controller 文件"""
    all_apis = []
    api_path = Path(base_path) / 'api'

    if not api_path.exists():
        print(f"❌ 错误: api 目录不存在")
        return all_apis

    for module_dir in api_path.iterdir():
        if not module_dir.is_dir():
            continue

        controller_dir = module_dir / 'controller'
        if not controller_dir.exists():
            continue

        for go_file in controller_dir.glob('*.go'):
            apis = extract_api_info_from_file(str(go_file))
            all_apis.extend(apis)

    return all_apis


def generate_go_file(apis: List[APIInfo], output_file: str = 'middleware/apiDescriptions.go'):
    """生成 Go 映射表文件"""

    # 按路由分组
    route_desc = {}
    route_methods = {}

    for api in apis:
        router = api.router
        method = api.method.upper()
        description = api.description or api.summary or ''

        if not router or not description:
            continue

        if router not in route_methods:
            route_methods[router] = {}
        route_methods[router][method] = description

    # 找出可以不区分方法的路由
    for router, methods in route_methods.items():
        if len(methods) == 1:
            route_desc[router] = list(methods.values())[0]
        else:
            descriptions = set(methods.values())
            if len(descriptions) == 1:
                route_desc[router] = list(descriptions)[0]

    # 找出需要区分方法的路由
    method_based = {}
    for router, methods in route_methods.items():
        if len(methods) > 1:
            descriptions = set(methods.values())
            if len(descriptions) > 1:
                for method, desc in methods.items():
                    key = f"{method.lower()}:{router}"
                    method_based[key] = desc

    # 收集动态路由
    dynamic_routes = []
    for api in apis:
        router = api.router
        description = api.description or api.summary or ''
        if router and description and '{' in router and '}' in router:
            pattern = router
            pattern = re.sub(r'\{id\}', r'\\d+', pattern)
            pattern = re.sub(r'\{[^}]+\}', r'[^/]+', pattern)
            pattern = f'^{pattern}$'
            dynamic_routes.append({'pattern': pattern, 'desc': description})

    # 生成 Go 代码
    content = '''// API描述映射
// 自动生成，请勿手动修改
// 使用 sync_api_descriptions.py 生成

package middleware

import (
\t"regexp"
\t"strings"
)

// GetAPIDescription 根据 URL 路径和 HTTP 方法获取 API 的中文描述
func GetAPIDescription(url, method string) string {
\t// 精确匹配
\tif desc, ok := getDescriptionMap()[url]; ok {
\t\treturn desc
\t}

\t// 方法+路径匹配
\tkey := method + ":" + url
\tif desc, ok := getMethodBasedDescriptionMap()[key]; ok {
\t\treturn desc
\t}

\t// 动态路由匹配
\treturn matchDynamicRoute(url)
}

// getDescriptionMap 返回 API 路径到描述的映射
func getDescriptionMap() map[string]string {
\tdescriptions := map[string]string{
'''

    # 添加路由描述
    for router in sorted(route_desc.keys()):
        desc = route_desc[router].replace('"', '\\"')
        content += f'\t\t"{router}": "{desc}",\n'

    content += '''\t}

\treturn descriptions
}

// getMethodBasedDescriptionMap 返回基于 HTTP Method + URL 的映射
func getMethodBasedDescriptionMap() map[string]string {
\treturn map[string]string{
'''

    # 添加方法+路径描述
    for key in sorted(method_based.keys()):
        desc = method_based[key].replace('"', '\\"')
        content += f'\t\t"{key}": "{desc}",\n'

    content += '''\t}
}

// matchDynamicRoute 匹配带参数的动态路由
func matchDynamicRoute(url string) string {
\tdynamicRoutes := []struct {
\t\tpattern *regexp.Regexp
\t\tdesc    string
\t}{
'''

    # 添加动态路由
    seen = set()
    for route in sorted(dynamic_routes, key=lambda x: x['pattern']):
        if route['pattern'] not in seen:
            seen.add(route['pattern'])
            desc = route['desc'].replace('"', '\\"')
            content += f'\t\t{{regexp.MustCompile(`{route["pattern"]}`), "{desc}"}},\n'

    content += '''\t}

\tfor _, route := range dynamicRoutes {
\t\tif route.pattern.MatchString(url) {
\t\t\treturn route.desc
\t\t}
\t}

\treturn extractActionFromURL(url)
}

// extractActionFromURL 从 URL 中提取操作类型
func extractActionFromURL(url string) string {
\tparts := strings.Split(url, "/")
\tif len(parts) >= 2 {
\t\taction := parts[len(parts)-1]
\t\tswitch action {
\t\tcase "add", "create":
\t\t\treturn "新增"
\t\tcase "update":
\t\t\treturn "修改"
\t\tcase "delete":
\t\t\treturn "删除"
\t\tcase "batchDelete", "batch-delete":
\t\t\treturn "批量删除"
\t\tcase "clean":
\t\t\treturn "清空"
\t\tcase "execute", "start":
\t\t\treturn "执行"
\t\tcase "stop":
\t\t\treturn "停止"
\t\tcase "deploy":
\t\t\treturn "部署"
\t\tcase "sync":
\t\t\treturn "同步"
\t\t}
\t}
\treturn ""
}
'''

    # 确保目录存在
    Path(output_file).parent.mkdir(parents=True, exist_ok=True)

    # 写入文件
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(content)


def print_statistics(apis: List[APIInfo]):
    """打印统计信息"""
    total = len(apis)
    with_desc = len([api for api in apis if api.description])
    without_desc = total - with_desc

    print(f"\n📊 统计信息:")
    print(f"  总 API 数: {total}")
    print(f"  有描述: {with_desc} ({with_desc*100//total}%)")
    print(f"  无描述: {without_desc} ({without_desc*100//total}%)")

    if without_desc > 0:
        print(f"\n⚠️  以下 API 缺少 @Description 注解 (前10个):")
        count = 0
        for api in apis:
            if not api.description:
                print(f"  - {api.method} {api.router}")
                print(f"    位置: {api.file_path}:{api.line_number}")
                count += 1
                if count >= 10:
                    if without_desc > 10:
                        print(f"  ... 还有 {without_desc - 10} 个")
                    break


def main():
    """主函数"""
    print("=" * 60)
    print("  API 描述自动同步工具")
    print("=" * 60)

    print("\n🔍 步骤 1/2: 扫描 API 文件...")
    apis = scan_controllers()

    if not apis:
        print("❌ 未找到任何 API 定义")
        return

    print(f"✓ 找到 {len(apis)} 个 API")

    print("\n📝 步骤 2/2: 更新 middleware/apiDescriptions.go...")
    generate_go_file(apis)
    print("✓ 已更新 middleware/apiDescriptions.go")

    print_statistics(apis)

    print("\n" + "=" * 60)
    print("  ✅ 完成!")
    print("=" * 60)
    print("\n下一步:")
    print("  1. 重启应用")
    print("  2. 访问任意 API")
    print("  3. 查看操作日志的 description 字段")
    print("")


if __name__ == '__main__':
    main()
