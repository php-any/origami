# 团队导航页系统

一个基于 Origami 的团队内部导航页系统，用于快速访问常用工具链接和项目不同环境的地址。

## 功能特性

- 🛠️ **常用工具导航**：快速访问 GitLab、Jira、Jenkins、Grafana 等常用工具
- 🌍 **项目环境管理**：集中管理项目在不同环境（开发、测试、预发布、生产）的访问地址
- 🔍 **快速搜索**：支持搜索工具和项目环境，支持键盘快捷键（Ctrl/Cmd + K）
- 📱 **响应式设计**：适配不同屏幕尺寸，支持移动端访问
- 🎨 **现代化 UI**：美观的暗色主题界面，提供良好的用户体验

## 快速开始

### 1. 安装依赖

```bash
cd examples/team-navigation
go mod download
```

### 2. 运行服务

```bash
go run main.go
```

服务将在 `http://127.0.0.1:8080` 启动

### 3. 访问页面

在浏览器中打开 `http://127.0.0.1:8080` 即可访问导航页

## 配置说明

### 常用工具配置

在 `pages/index.html` 中的 `$commonTools` 数组中配置常用工具链接：

```zy
$commonTools = [
    { name: "GitLab", url: "https://gitlab.example.com", icon: "🔗", category: "代码管理", desc: "代码仓库和 CI/CD" },
    // 添加更多工具...
];
```

### 项目环境配置

在 `pages/index.html` 中的 `$projectEnvs` 数组中配置项目环境链接：

```zy
$projectEnvs = [
    {
        project: "Origami",
        envs: [
            { env: "开发环境", url: "http://dev-origami.example.com", status: "运行中", color: "green" },
            // 添加更多环境...
        ]
    },
    // 添加更多项目...
];
```

环境状态支持的颜色：

- `green`：运行中
- `yellow`：维护中
- `red`：异常

## 项目结构

```
team-navigation/
├── main.go              # Go 程序入口
├── http.zy              # HTTP 服务器配置
├── main.zy              # 请求处理逻辑
├── pages/               # 页面文件目录
│   ├── index.html       # 导航页主页面
│   └── assets/          # 静态资源
│       ├── css/
│       │   └── styles.css    # 样式文件
│       └── js/
│           └── app.js         # 交互脚本
├── go.mod               # Go 模块配置
└── README.md            # 说明文档
```

## 自定义修改

### 修改端口

在 `http.zy` 中修改端口号：

```zy
$server = new Server("0.0.0.0", port: 8080);
```

### 修改样式

编辑 `pages/assets/css/styles.css` 文件，可以自定义颜色、字体、布局等样式。

### 添加新功能

- 页面逻辑：修改 `pages/index.html` 中的脚本部分
- 交互功能：修改 `pages/assets/js/app.js`
- 后端处理：修改 `main.zy` 中的请求处理逻辑

## 键盘快捷键

- `Ctrl/Cmd + K`：快速聚焦搜索框

## 技术栈

- **后端**：Origami (基于 Go 的脚本语言)
- **前端**：HTML5 + CSS3 + JavaScript
- **服务器**：Origami HTTP 标准库

## 参考示例

本示例参考了：

- `examples/html`：基础的 HTML 页面服务示例
- `examples/spring`：注解驱动的路由示例

## 许可证

本示例遵循项目主许可证。
