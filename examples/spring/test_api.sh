# Spring 示例 API 测试脚本
# 
# 使用方法：
# 1. 首先运行服务器: ./origami run examples/spring/index.php
# 2. 在另一个终端运行此脚本进行测试

# 基础 URL
BASE_URL="http://127.0.0.1:8080"

echo "======================================"
echo "Spring 示例 API 测试"
echo "======================================"
echo ""

# 1. 测试 Hello 接口
echo "1. 测试 GET /api/hello"
curl -s -X GET "$BASE_URL/api/hello" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 2. 测试 Info 接口
echo "2. 测试 GET /api/info"
curl -s -X GET "$BASE_URL/api/info" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 3. 测试 Status 接口
echo "3. 测试 GET /api/status"
curl -s -X GET "$BASE_URL/api/status" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 4. 获取用户列表
echo "4. 测试 GET /api/users"
curl -s -X GET "$BASE_URL/api/users" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 5. 获取单个用户
echo "5. 测试 GET /api/user/1"
curl -s -X GET "$BASE_URL/api/user/1" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 6. 获取商品列表
echo "6. 测试 GET /api/products"
curl -s -X GET "$BASE_URL/api/products" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 7. 获取单个商品
echo "7. 测试 GET /api/product/1"
curl -s -X GET "$BASE_URL/api/product/1" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 8. 创建新商品
echo "8. 测试 POST /api/products"
curl -s -X POST "$BASE_URL/api/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试商品",
    "price": 999.99,
    "category": "测试分类",
    "description": "这是一个测试商品"
  }' | python3 -m json.tool
echo ""
echo "---"
echo ""

# 9. 更新商品
echo "9. 测试 PUT /api/product/1"
curl -s -X PUT "$BASE_URL/api/product/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro Max",
    "price": 8999.00
  }' | python3 -m json.tool
echo ""
echo "---"
echo ""

# 10. 搜索商品
echo "10. 测试 GET /api/products/search?keyword=iPhone"
curl -s -X GET "$BASE_URL/api/products/search?keyword=iPhone" | python3 -m json.tool
echo ""
echo "---"
echo ""

# 11. 用户登录
echo "11. 测试 POST /api/auth/login"
curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }' | python3 -m json.tool
echo ""
echo "---"
echo ""

# 12. 用户注册
echo "12. 测试 POST /api/auth/register"
curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "password123",
    "email": "newuser@example.com"
  }' | python3 -m json.tool
echo ""
echo "---"
echo ""

echo "======================================"
echo "测试完成！"
echo "======================================"
