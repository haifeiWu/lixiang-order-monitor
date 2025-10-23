#!/bin/bash
# Bark 推送测试脚本

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Bark 推送测试${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 从配置文件读取 Bark URL
CONFIG_FILE="../../config.yaml"

if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}❌ 配置文件不存在: $CONFIG_FILE${NC}"
    echo -e "${YELLOW}请先创建配置文件${NC}"
    exit 1
fi

# 读取 Bark 配置
BARK_URL=$(grep "bark_server_url:" "$CONFIG_FILE" | sed 's/.*: *"\?\([^"]*\)"\?.*/\1/' | tr -d ' ')
BARK_SOUND=$(grep "bark_sound:" "$CONFIG_FILE" | sed 's/.*: *"\?\([^"]*\)"\?.*/\1/' | tr -d ' ')
BARK_GROUP=$(grep "bark_group:" "$CONFIG_FILE" | sed 's/.*: *"\?\([^"]*\)"\?.*/\1/' | tr -d ' ')

if [ -z "$BARK_URL" ]; then
    echo -e "${RED}❌ Bark URL 未配置${NC}"
    echo -e "${YELLOW}请在 config.yaml 中配置 bark_server_url${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 找到 Bark 配置${NC}"
echo -e "URL: ${BARK_URL}"
echo -e "Sound: ${BARK_SOUND:-minuet}"
echo -e "Group: ${BARK_GROUP:-lixiang-monitor}\n"

# 测试推送
echo -e "${YELLOW}📱 发送测试推送...${NC}\n"

RESPONSE=$(curl --location "$BARK_URL" \
--silent \
--write-out "\nHTTP_STATUS:%{http_code}" \
--header 'Content-Type: application/json; charset=utf-8' \
--data "{
  \"title\": \"理想汽车订单监控\",
  \"body\": \"Bark 推送测试成功！\\n\\n测试时间: $(date '+%Y-%m-%d %H:%M:%S')\",
  \"sound\": \"${BARK_SOUND:-minuet}\",
  \"group\": \"${BARK_GROUP:-lixiang-monitor}\"
}")

# 分离响应体和状态码
HTTP_BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS:.*//')
HTTP_STATUS=$(echo "$RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

echo -e "${BLUE}响应状态码: ${HTTP_STATUS}${NC}"
echo -e "${BLUE}响应内容:${NC}"
echo "$HTTP_BODY" | python3 -m json.tool 2>/dev/null || echo "$HTTP_BODY"
echo ""

# 判断结果
if [ "$HTTP_STATUS" = "200" ]; then
    echo -e "${GREEN}✅ Bark 推送测试成功！${NC}"
    echo -e "${GREEN}请检查你的 iPhone/iPad/Mac 是否收到推送${NC}"
else
    echo -e "${RED}❌ Bark 推送测试失败${NC}"
    echo -e "${YELLOW}请检查:${NC}"
    echo -e "  1. Bark Server URL 是否正确"
    echo -e "  2. Bark 服务器是否可访问"
    echo -e "  3. Bark App 是否已添加该服务器"
    exit 1
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}测试完成${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${YELLOW}提示:${NC}"
echo "• 如果未收到推送，请检查 Bark App 的通知权限"
echo "• 确认 Bark App 中已添加对应的服务器"
echo "• 可以在 Bark App 中查看推送历史"
echo ""
