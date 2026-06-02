# Spring 示例 API 测试脚本 (PowerShell)
# 
# 使用方法：
# 1. 首先运行服务器: .\origami.exe run examples\spring\index.php
# 2. 在另一个 PowerShell 窗口运行此脚本进行测试

# 基础 URL
$BASE_URL = "http://127.0.0.1:8080"

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Spring 示例 API 测试" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# 辅助函数：发送请求并格式化 JSON
function Test-Api {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Description,
        [hashtable]$Body = $null
    )
    
    Write-Host "$Description" -ForegroundColor Yellow
    
    try {
        $params = @{
            Method = $Method
            Uri = $Url
            ContentType = "application/json; charset=utf-8"
        }
        
        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json -Depth 10)
        }
        
        $response = Invoke-RestMethod @params
        $response | ConvertTo-Json -Depth 10
    } catch {
        Write-Host "错误: $_" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "---" -ForegroundColor Gray
    Write-Host ""
}

# 1. 测试 Hello 接口
Test-Api -Method "GET" -Url "$BASE_URL/api/hello" -Description "1. 测试 GET /api/hello"

# 2. 测试 Info 接口
Test-Api -Method "GET" -Url "$BASE_URL/api/info" -Description "2. 测试 GET /api/info"

# 3. 测试 Status 接口
Test-Api -Method "GET" -Url "$BASE_URL/api/status" -Description "3. 测试 GET /api/status"

# 4. 获取用户列表
Test-Api -Method "GET" -Url "$BASE_URL/api/users" -Description "4. 测试 GET /api/users"

# 5. 获取单个用户
Test-Api -Method "GET" -Url "$BASE_URL/api/user/1" -Description "5. 测试 GET /api/user/1"

# 6. 获取商品列表
Test-Api -Method "GET" -Url "$BASE_URL/api/products" -Description "6. 测试 GET /api/products"

# 7. 获取单个商品
Test-Api -Method "GET" -Url "$BASE_URL/api/product/1" -Description "7. 测试 GET /api/product/1"

# 8. 创建新商品
Test-Api -Method "POST" -Url "$BASE_URL/api/products" -Description "8. 测试 POST /api/products" -Body @{
    name = "测试商品"
    price = 999.99
    category = "测试分类"
    description = "这是一个测试商品"
}

# 9. 更新商品
Test-Api -Method "PUT" -Url "$BASE_URL/api/product/1" -Description "9. 测试 PUT /api/product/1" -Body @{
    name = "iPhone 15 Pro Max"
    price = 8999.00
}

# 10. 搜索商品
Test-Api -Method "GET" -Url "$BASE_URL/api/products/search?keyword=iPhone" -Description "10. 测试 GET /api/products/search?keyword=iPhone"

# 11. 用户登录
Test-Api -Method "POST" -Url "$BASE_URL/api/auth/login" -Description "11. 测试 POST /api/auth/login" -Body @{
    username = "admin"
    password = "123456"
}

# 12. 用户注册
Test-Api -Method "POST" -Url "$BASE_URL/api/auth/register" -Description "12. 测试 POST /api/auth/register" -Body @{
    username = "newuser"
    password = "password123"
    email = "newuser@example.com"
}

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "测试完成！" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Cyan
