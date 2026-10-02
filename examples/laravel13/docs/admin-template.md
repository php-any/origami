# 后台前端模板

后台使用固定版本 `@tabler/core 1.6.1`（MIT）。整体导航、登录卡片、工作台网格、统计卡片、表格、进度条、空状态、常用入口均使用官方组件与样式。商品等业务表单和资源列表继续使用 Filament 5 / Livewire 4，沿用现有认证、Policy、校验、上传和业务操作。

官方参考：

- https://docs.tabler.io/ui/getting-started/frameworks/laravel
- https://docs.tabler.io/ui/layout/page-layouts
- https://github.com/tabler/tabler

## 文件位置

- `resources/views/vendor/filament-panels/components/layout/`：应用层 Blade 布局覆盖，不修改 Composer vendor。
- `resources/views/vendor/filament-panels/livewire/topbar.blade.php`：官方顶部导航结构，保留 Filament 全局搜索、通知、个人资料与退出。
- `resources/views/filament/pages/dashboard.blade.php`：Tabler 原生组件接入真实、经过权限过滤的业务数据。
- `resources/css/admin/tabler.css`：Vite 构建入口，由 `vite.config.js` 从安装包读取官方完整 CSS。
- `resources/css/filament/admin/theme.css`：官方 Filament CSS、浅色后台导航外壳，以及仅作用于系统设置页的分行配置面板与底部操作区。
- `resources/js/admin.js`：加载 Tabler 官方交互；同步官方折叠侧栏偏好；短屏桌面导航自动滚动到当前模块。

## 样式边界

Vite 将 Tabler 官方 CSS 包装在 `@scope (.tabler-ui) to (.filament-content)` 内，并将 HTML / body / root 选择器映射到作用域根。使用当前支持 CSS scope 的浏览器；Tabler 1.6 本身也使用 `light-dark()`、OKLCH 等现代 CSS。

`.filament-content` 为 Filament 表单、资源表格、搜索、通知和弹窗的边界。工作台在边界内建立独立 `.tabler-ui` 作用域，使用 Tabler 组件。不要把 Bootstrap 的重置样式扩散到业务组件，或重新添加大段自定义主题。

导航由 `filament()->getNavigation()` 生成，遵守现有权限；数据来自数据库。空数据保留官方空状态，不添加演示订单或虚构增长率。

顶部仅提供工作台入口与全局操作，不渲染固定层级。面包屑统一由 Filament 页面标题区生成，沿用资源页面的真实链接与当前页面标签。

## 构建与验收

```sh
pnpm install --frozen-lockfile
pnpm run build
go run -mod=mod . serve --port=18086
go run -mod=mod . run tests/origami/workbench_access.php
```

应用层 PHP / Blade 更新后重启开发中的 Origami 服务，使进程内解析缓存刷新。`go run -mod=mod . view:clear` 可清除 Laravel 编译视图。
