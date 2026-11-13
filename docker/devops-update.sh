#!/usr/bin/env bash

# Update only the selected AutoOps frontend/backend containers.
# Usage: ./devops-update.sh <version> <web|api> [web|api]

set -Eeuo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

REGISTRY="crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s"
API_REPOSITORY="${AUTOOPS_API_REPOSITORY:-${REGISTRY}/deviops-api}"
WEB_REPOSITORY="${AUTOOPS_WEB_REPOSITORY:-${REGISTRY}/deviops-web}"
WAIT_TIMEOUT="${AUTOOPS_WAIT_TIMEOUT:-180}"
PODMAN_DOCKER_NOTICE="Emulate Docker CLI using podman. Create /etc/containers/nodocker to quiet msg."

docker_cli() {
    docker "$@" 2> >(
        while IFS= read -r line || [ -n "$line" ]; do
            [ "$line" = "$PODMAN_DOCKER_NOTICE" ] || printf '%s\n' "$line" >&2
        done
    )
}

usage() {
    cat <<EOF
用法: $0 <version> <web|api> [web|api]

示例:
  $0 v6.0 web          只更新 Web
  $0 v6.0 api          只更新 API
  $0 v6.0 web api      同时更新 Web 和 API

只更新指定镜像和容器，不会重启 MySQL、Redis、Kafka、Vector 或监控服务。
EOF
}

fail() {
    echo -e "${RED}错误: $*${NC}" >&2
    exit 1
}

on_error() {
    local exit_code=$?
    echo -e "${RED}更新失败，出错行: ${BASH_LINENO[0]}，退出码: ${exit_code}${NC}" >&2
    exit "$exit_code"
}
trap on_error ERR

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
    usage
    exit 1
fi

VERSION="$1"
shift
[[ "$VERSION" =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]] || fail "镜像版本格式无效: $VERSION"
[[ "$WAIT_TIMEOUT" =~ ^[0-9]+$ ]] || fail "AUTOOPS_WAIT_TIMEOUT 必须是正整数"

UPDATE_WEB=false
UPDATE_API=false
for component in "$@"; do
    case "$component" in
        web)
            [ "$UPDATE_WEB" = false ] || fail "web 参数重复"
            UPDATE_WEB=true
            ;;
        api)
            [ "$UPDATE_API" = false ] || fail "api 参数重复"
            UPDATE_API=true
            ;;
        *)
            usage
            fail "未知组件: $component（只能是 web 或 api）"
            ;;
    esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"
ENV_FILE="$SCRIPT_DIR/.env"
[ -f "$COMPOSE_FILE" ] || fail "找不到 $COMPOSE_FILE，请先在 docker 目录初始化部署"
[ -f "$ENV_FILE" ] || fail "找不到 $ENV_FILE，请先运行 devops-start.sh"
command -v docker >/dev/null 2>&1 || fail "未安装 Docker 或 Podman"
docker_cli info >/dev/null 2>&1 || fail "Docker/Podman 未运行或当前用户无访问权限"

if docker_cli compose version >/dev/null 2>&1; then
    COMPOSE=(docker_cli compose)
    COMPOSE_DISPLAY="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE=(docker-compose)
    COMPOSE_DISPLAY="docker-compose"
else
    fail "找不到 docker compose 或 docker-compose"
fi

get_env_value() {
    local key="$1"
    awk -v key="$key" 'index($0, key "=") == 1 { print substr($0, length(key) + 2); exit }' "$ENV_FILE"
}

is_valid_port() {
    local port="$1"
    [[ "$port" =~ ^[0-9]{1,5}$ ]] || return 1
    local numeric=$((10#$port))
    (( numeric >= 1 && numeric <= 65535 ))
}

resolve_web_port() {
    local port
    port="$(docker_cli inspect --format '{{with index .HostConfig.PortBindings "80/tcp"}}{{(index . 0).HostPort}}{{end}}' autoops-web 2>/dev/null || true)"
    if is_valid_port "$port"; then
        printf '%s' "$port"
        return
    fi

    port="$(get_env_value WEB_PORT)"
    port="${port:-8088}"
    is_valid_port "$port" || fail "Web 端口无效: $port"
    printf '%s' "$port"
}

set_env_value() {
    local key="$1"
    local value="$2"
    local temp_file
    temp_file="$(mktemp "$SCRIPT_DIR/.env.XXXXXX")"
    ENV_KEY="$key" ENV_VALUE="$value" awk '
        BEGIN { key = ENVIRON["ENV_KEY"]; value = ENVIRON["ENV_VALUE"]; found = 0 }
        index($0, key "=") == 1 {
            if (!found) { print key "=" value; found = 1 }
            next
        }
        { print }
        END { if (!found) print key "=" value }
    ' "$ENV_FILE" > "$temp_file"
    chmod 600 "$temp_file"
    mv "$temp_file" "$ENV_FILE"
}

WEB_PORT=""
if [ "$UPDATE_WEB" = true ]; then
    WEB_PORT="$(resolve_web_port)"
    export WEB_PORT
fi

declare -a IMAGES=()
declare -a SERVICES=()
declare -a CONTAINERS=()
if [ "$UPDATE_WEB" = true ]; then
    WEB_IMAGE="${WEB_REPOSITORY}:${VERSION}"
    IMAGES+=("$WEB_IMAGE")
    SERVICES+=(autoops-web)
    CONTAINERS+=(autoops-web)
fi
if [ "$UPDATE_API" = true ]; then
    API_IMAGE="${API_REPOSITORY}:${VERSION}"
    IMAGES+=("$API_IMAGE")
    SERVICES+=(autoops-api)
    CONTAINERS+=(autoops-api)
fi

cd "$SCRIPT_DIR"
echo -e "${GREEN}AutoOps 选择性更新${NC}"
printf '版本: %s\n组件: %s\n' "$VERSION" "${SERVICES[*]}"
if [ "$UPDATE_WEB" = true ]; then
    printf 'Web 端口: %s\n' "$WEB_PORT"
fi

"${COMPOSE[@]}" config >/dev/null

echo -e "${YELLOW}[1/3] 拉取目标镜像...${NC}"
for image in "${IMAGES[@]}"; do
    docker_cli pull "$image"
done

echo -e "${YELLOW}[2/3] 更新镜像版本配置...${NC}"
if [ "$UPDATE_WEB" = true ]; then
    set_env_value AUTOOPS_WEB_IMAGE "$WEB_IMAGE"
    set_env_value WEB_PORT "$WEB_PORT"
fi
if [ "$UPDATE_API" = true ]; then
    set_env_value AUTOOPS_API_IMAGE "$API_IMAGE"
fi
"${COMPOSE[@]}" config >/dev/null

echo -e "${YELLOW}[3/3] 只重建目标服务...${NC}"
"${COMPOSE[@]}" up -d --no-deps --force-recreate "${SERVICES[@]}"

for container in "${CONTAINERS[@]}"; do
    deadline=$((SECONDS + WAIT_TIMEOUT))
    ready=false
    while (( SECONDS < deadline )); do
        if ! docker_cli inspect "$container" >/dev/null 2>&1; then
            sleep 2
            continue
        fi

        status="$(docker_cli inspect --format '{{.State.Status}}' "$container")"
        health="$(docker_cli inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$container")"
        if [ "$status" = running ] && { [ "$health" = none ] || [ "$health" = healthy ]; }; then
            echo -e "${GREEN}✓ $container 已就绪（health=$health）${NC}"
            ready=true
            break
        fi
        if [ "$status" = exited ] || [ "$status" = dead ] || [ "$health" = unhealthy ]; then
            docker_cli logs --tail 80 "$container" 2>&1 || true
            fail "$container 启动失败（status=$status, health=$health）"
        fi
        sleep 2
    done
    if [ "$ready" != true ]; then
        docker_cli logs --tail 80 "$container" 2>&1 || true
        fail "等待 $container 就绪超时"
    fi
done

echo -e "${GREEN}✓ 更新完成。其他服务未重启。${NC}"
echo "查看容器状态: $COMPOSE_DISPLAY ps ${SERVICES[*]}"
echo "查看容器日志: $COMPOSE_DISPLAY logs -f ${SERVICES[*]}"
