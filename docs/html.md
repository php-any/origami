# HTML 渲染与模板

折言 Origami 内置 HTML 模板解析能力，可直接渲染 `.html` 文件中的模板语法与内嵌脚本，并通过标准库 `Net\Http` 快速输出到浏览器。

- 渲染入口：`$response->view($templatePath, $data)`
- 应用启动：`Net\Http\app($request, $response, filePath: "./src/main.zy", fun: "App\\main")`
- 静态资源：`$server->static(prefix: "/assets/", dir: "./pages")`

## 快速上手

```zy
use Net\Http\Server;
use Net\Http\app;

$server = new Server("0.0.0.0", port: 8080);

// 静态资源（可选）：将 /assets/ 映射到 ./pages 目录
$server->static(prefix: "/assets/", dir: "./pages");

// 统一入口，委托到 main.zy 中的 App\\main 函数
$server->any(($request, $response) => {
    app($request, $response, "./main.zy", "App\\main")
});

$server->run();
```

`main.zy` 中定义应用入口函数，并渲染页面：

```zy
namespace App;

function main($request, $response) {
    // 传入模板上下文对象（将以变量形式在模板中可用）
    $ctx = {
        title: "首页",
        features: [
            { title: "极速", desc: "毫秒级响应" },
            { title: "易用", desc: "语法直观" }
        ]
    };
    $response->view("./pages/index.html", $ctx);
}
```

## 模板语法

HTML 文件中可直接使用以下能力：

- 表达式插值：`{$expr}`
- 控制指令（以 HTML 属性形式存在）：
  - `for`：循环块。例如：`<div for="$f in $features">{$f->title}</div>`
  - `if` / `elseif` / `else`：条件渲染。例如：`<p if="$user != null">...</p>`
 - 动态属性（:ANY）：任意属性“名称”以冒号前缀表示“表达式求值”，如 `:title="$user->name"`、`:data-count="count($list)"`、`:style="$active ? 'color:red' : ''"`。
- 内嵌脚本块（可选）：
  - 在 `<script type="text/zy"> ... </script>` 中编写折言脚本，用于准备页面局部数据。

一个最小模板示例：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <title>{$title}</title>
  <link rel="stylesheet" href="/assets/css/styles.css" />
  <script type="text/zy">
    // 可在此定义页面私有数据或计算
    $greeting = "欢迎使用 Origami";
  </script>
  </head>
<body>
  <h1>{$greeting}</h1>
  <section>
    <div for="$f in $features">
      <h3>{$f->title}</h3>
      <p>{$f->desc}</p>
    </div>
  </section>
  <script src="/assets/js/app.js"></script>
</body>
</html>
```

### 动态属性（:ANY）示例

以下示例展示对任意属性使用 `:表达式` 动态求值：

```html
<a href="/profile" :title="$user ? '用户：' + $user->name : '游客'">
  <span class="badge" :data-count="count($notifications)">通知</span>
  <i class="dot" :style="$online ? 'background:#22c55e' : 'background:#9ca3af'"></i>
</a>
```

说明：
- 以冒号开头的“属性名”会被解析为折言表达式并求值，其结果再写入对应的实际属性名。
- 适用于任意属性名（class、style、data-*、aria-*、自定义属性等）。

## 路由与静态资源

- 注册统一路由：`$server->any(($req,$res)=>{ app($req,$res,"./main.zy","App\\main") });`
- 静态资源：`$server->static("/assets/", "./pages")`，模板内即可引用 `/assets/...` 路径。

> 完整示例可参考 `examples/html/` 目录。

## 浏览器端演示

无需本地运行，也可以使用在线演示工具体验语法与渲染流程：

- 在线演示：`https://php-any.github.io/wasm-demo/`

该演示支持在浏览器中编写并运行折言脚本，快速验证语法与输出。


