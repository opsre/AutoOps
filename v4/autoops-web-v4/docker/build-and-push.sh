#!/bin/bash

# DevOps Web 镜像构建和推送脚本
# Author: zhangfan
# Date: 2025-10-15

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置变量
VERSION="v4.0"  # 镜像版本号
IMAGE_REGISTRY="crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com"
IMAGE_NAMESPACE="zhangfan_k8s"
IMAGE_NAME="deviops-web"
IMAGE_TAG="${1:-${VERSION}}"  # 优先使用参数传入的标签，否则使用 VERSION
FULL_IMAGE_NAME="${IMAGE_REGISTRY}/${IMAGE_NAMESPACE}/${IMAGE_NAME}:${IMAGE_TAG}"

# 项目根目录（脚本所在目录的上一级）
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST_DIR="${PROJECT_ROOT}/dist"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  DevOps Web 镜像构建和推送工具${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${YELLOW}镜像名称: ${FULL_IMAGE_NAME}${NC}"
echo -e "${YELLOW}项目路径: ${PROJECT_ROOT}${NC}"
echo ""

# 步骤 1: 构建前端项目
echo -e "${GREEN}[1/5] 开始构建前端项目...${NC}"
cd "${PROJECT_ROOT}"
npm run build

if [ ! -d "${DIST_DIR}" ]; then
    echo -e "${RED}错误: dist 目录不存在，构建失败！${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 前端项目构建完成${NC}"
echo ""

# 步骤 2: 构建 Docker 镜像
echo -e "${GREEN}[2/5] 开始构建 Docker 镜像...${NC}"
docker build -t "${FULL_IMAGE_NAME}" -f "${PROJECT_ROOT}/docker/Dockerfile" "${PROJECT_ROOT}"

if [ $? -ne 0 ]; then
    echo -e "${RED}错误: Docker 镜像构建失败！${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker 镜像构建完成${NC}"
echo ""

# 步骤 3: 推送镜像到远程仓库
echo -e "${GREEN}[3/5] 开始推送镜像到远程仓库...${NC}"
docker push "${FULL_IMAGE_NAME}"

if [ $? -ne 0 ]; then
    echo -e "${RED}错误: 镜像推送失败！${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 镜像推送完成${NC}"
echo ""

# 步骤 4: 删除本地 dist 目录
echo -e "${GREEN}[4/5] 清理本地 dist 目录...${NC}"
if [ -d "${DIST_DIR}" ]; then
    rm -rf "${DIST_DIR}"
    echo -e "${GREEN}✓ dist 目录已删除${NC}"
else
    echo -e "${YELLOW}! dist 目录不存在，跳过删除${NC}"
fi
echo ""

# 步骤 5: 删除本地 Docker 镜像
echo -e "${GREEN}[5/5] 删除本地 Docker 镜像...${NC}"
docker rmi "${FULL_IMAGE_NAME}"

if [ $? -ne 0 ]; then
    echo -e "${YELLOW}! 本地镜像删除失败（可能镜像不存在或正在使用）${NC}"
else
    echo -e "${GREEN}✓ 本地镜像已删除${NC}"
fi
echo ""

# 完成
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  构建和推送完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${YELLOW}镜像地址: ${FULL_IMAGE_NAME}${NC}"
echo ""
echo -e "${YELLOW}使用方式:${NC}"
echo -e "  docker pull ${FULL_IMAGE_NAME}"
echo -e "  docker run -d -p 80:80 ${FULL_IMAGE_NAME}"
echo ""
