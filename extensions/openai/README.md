# OpenAI 标准库

OpenAI 标准库为折言提供了调用 OpenAI API 的能力，包括 Chat Completions、文本嵌入、图片生成、语音合成和语音转文字等功能。

## 概述

OpenAI 模块通过封装 [openai-go](https://github.com/openai/openai-go) SDK 实现，注册为 `OpenAI\Client` 类。使用前需要有效的 OpenAI API Key。

**提供的类：**

| 类名 | 说明 |
|------|------|
| `OpenAI\Client` | 主客户端类，所有 API 调用的入口 |
| `OpenAI\ChatCompletion` | Chat Completions 响应结果 |
| `OpenAI\EmbeddingResult` | 文本嵌入响应结果 |
| `OpenAI\ImageResult` | 图片生成响应结果 |
| `OpenAI\TranscriptionResult` | 语音转文字响应结果 |

## openai-go SDK API 映射

本扩展封装 [openai-go](https://github.com/openai/openai-go) SDK，提供与 SDK README 示例对齐的脚本层 API。

### Chat Completions 参数映射

| openai-go SDK 参数 | 折言 opts 参数 | 类型 | 说明 |
|---|---|---|---|
| `Temperature` | `temperature` | float | 0~2，越高越随机 |
| `MaxTokens` | `maxTokens` | int | 最大输出 token 数 |
| `TopP` | `topP` | float | 核采样参数 |
| `Stop` | `stop` | array | 停止词列表 |
| `Seed` | `seed` | int | 确定性采样种子 |
| `FrequencyPenalty` | `frequencyPenalty` | float | -2.0~2.0，频率惩罚 |
| `PresencePenalty` | `presencePenalty` | float | -2.0~2.0，存在惩罚 |
| `MaxCompletionTokens` | `maxCompletionTokens` | int | 完成最大 token 数（新版） |
| `N` | `n` | int | 生成几个候选回复 |
| `ResponseFormat` | `responseFormat` | string/array | `"json_object"` 或 json_schema 对象 |

### 消息角色映射

| openai-go SDK | 折言 role 参数 |
|---|---|
| `openai.SystemMessage(content)` | `["role" => "system", "content" => "..."]` |
| `openai.UserMessage(content)` | `["role" => "user", "content" => "..."]` |
| `openai.AssistantMessage(content)` | `["role" => "assistant", "content" => "..."]` |
| `openai.DeveloperMessage(content)` | `["role" => "developer", "content" => "..."]` |
| `openai.ToolMessage(toolCallID, content)` | `["role" => "tool", "content" => "...", "tool_call_id" => "..."]` |

### API 方法映射

| openai-go SDK | 折言方法 |
|---|---|
| `client.Chat.Completions.New(ctx, params)` | `$client->chat($model, $messages, $opts)` |
| `client.Embeddings.New(ctx, params)` | `$client->embeddings($model, $input, $opts)` |
| `client.Images.Generate(ctx, params)` | `$client->images($model, $prompt, $opts)` |
| `client.Audio.Speech.New(ctx, params)` | `$client->speech($model, $input, $voice, $output, $opts)` |
| `client.Audio.Transcriptions.New(ctx, params)` | `$client->transcription($model, $file, $opts)` |

## 基本用法

### 创建客户端

```php
<?php
// 使用环境变量 OPENAI_API_KEY
$client = new OpenAI\Client("sk-your-api-key");

// 指定自定义 API 地址（如使用代理或兼容服务）
$client = new OpenAI\Client("sk-your-api-key", "https://api.example.com/v1");
```

### Chat Completions

发送对话请求并获取回复：

```php
<?php
$client = new OpenAI\Client(getenv("OPENAI_API_KEY"));

$messages = [
    ["role" => "system", "content" => "你是一个有帮助的助手。"],
    ["role" => "user", "content" => "用一句话介绍折言语言。"],
];

$result = $client->chat("gpt-4o", $messages);

echo $result->content;           // 回复内容
echo $result->role;              // 角色（始终为 "assistant"）
echo $result->finishReason;      // 完成原因（如 "stop"）
echo $result->usage->totalTokens; // 总 token 消耗
```

**可选参数（第三个参数）：**

```php
<?php
$options = [
    "temperature" => 0.7,
    "maxTokens" => 1000,
];

$result = $client->chat("gpt-4o", $messages, $options);
```

### 文本嵌入（Embeddings）

将文本转换为向量表示：

```php
<?php
$client = new OpenAI\Client(getenv("OPENAI_API_KEY"));

$result = $client->embeddings("text-embedding-3-small", "Hello world");

echo $result->dimensions;        // 向量维度
echo $result->usage->totalTokens;

// 获取向量数组
$vector = $result->embedding;    // float 数组
```

支持传入字符串数组批量处理：

```php
<?php
$result = $client->embeddings("text-embedding-3-small", ["Hello", "World"]);
```

### 图片生成

根据文字描述生成图片：

```php
<?php
$client = new OpenAI\Client(getenv("OPENAI_API_KEY"));

$result = $client->images("dall-e-3", "一只戴帽子的猫");

echo $result->url;              // 图片 URL
echo $result->revisedPrompt;    // 修订后的提示词
```

**可选参数：**

```php
<?php
$options = [
    "size" => "1024x1024",    // 可选: 256x256, 512x512, 1024x1024, 1536x1024, 1024x1536
    "quality" => "standard",  // 可选: standard, hd
];

$result = $client->images("dall-e-3", "一只戴帽子的猫", $options);
```

### 语音合成（Text-to-Speech）

将文字转换为语音文件：

```php
<?php
$client = new OpenAI\Client(getenv("OPENAI_API_KEY"));

$client->speech("tts-1", "你好，世界！", "alloy", "output.mp3");
echo "语音已保存到 output.mp3\n";
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| model | string | 模型名，如 `tts-1`、`tts-1-hd`、`gpt-4o-mini-tts` |
| input | string | 要转换的文字 |
| voice | string | 声音类型：alloy, ash, ballad, coral, echo, fable, onyx, nova, sage, shimmer, verse, marin, cedar |
| output | string | 输出文件路径 |

**可选参数：**

```php
<?php
$options = ["speed" => 1.0];  // 语速，0.25 到 4.0
$client->speech("tts-1", "你好", "alloy", "output.mp3", $options);
```

### 语音转文字（Transcription）

将音频文件转换为文字：

```php
<?php
$client = new OpenAI\Client(getenv("OPENAI_API_KEY"));

$result = $client->transcription("whisper-1", "audio.mp3");

echo $result->text;       // 转录文本
echo $result->language;   // 语言
echo $result->duration;   // 音频时长（秒）
```

**可选参数：**

```php
<?php
$options = ["language" => "zh"];  // 指定语言
$result = $client->transcription("whisper-1", "audio.mp3", $options);
```

## 类型参考

### OpenAI\ChatCompletion

| 属性 | 类型 | 说明 |
|------|------|------|
| content | string | 回复文本内容 |
| role | string | 消息角色（始终为 "assistant"） |
| finishReason | string | 完成原因（"stop"、"length" 等） |
| usage | object | token 使用量 |
| usage.promptTokens | int | 输入 token 数 |
| usage.completionTokens | int | 输出 token 数 |
| usage.totalTokens | int | 总 token 数 |

### OpenAI\EmbeddingResult

| 属性 | 类型 | 说明 |
|------|------|------|
| embedding | array | 向量数组（float[]） |
| dimensions | int | 向量维度 |
| usage | object | token 使用量 |

### OpenAI\ImageResult

| 属性 | 类型 | 说明 |
|------|------|------|
| url | string | 生成图片的 URL |
| b64Json | string | Base64 编码的图片数据（如适用） |
| revisedPrompt | string | 修订后的提示词 |

### OpenAI\TranscriptionResult

| 属性 | 类型 | 说明 |
|------|------|------|
| text | string | 转录文本 |
| language | string | 识别的语言 |
| duration | float | 音频时长（秒） |

## 错误处理

所有方法在 API 调用失败时会抛出异常，建议使用 try/catch 捕获：

```php
<?php
try {
    $client = new OpenAI\Client("invalid-key");
    $result = $client->chat("gpt-4o", [["role" => "user", "content" => "Hi"]]);
} catch (\Exception $e) {
    echo "错误: " . $e->getMessage() . "\n";
}
```

## 安装使用

OpenAI 模块作为**独立扩展**管理，不会随折言标准库一起编译。使用时需要将其作为 Go module 引入。

### 方法：创建自定义入口并注册扩展

1. 在你的项目中创建 `main.go`，导入并注册 OpenAI 扩展：

```go
package main

import (
    "os"

    _ "github.com/go-sql-driver/mysql"
    _ "modernc.org/sqlite"

    "github.com/php-any/origami/cmd"
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/std"
    "github.com/php-any/origami/std/php"
    "github.com/php-any/origami/std/system"
    openai "github.com/php-any/origami-openai"
)

func init() {
    cmd.SetRuntimeLoader(func(vm data.VM) {
        std.Load(vm)
        php.Load(vm)
        system.Load(vm)
        openai.Load(vm)  // 注册 OpenAI\Client 类
    })
}

func main() {
    if len(os.Args) > 1 && cmd.IsDirectScriptArg(os.Args[1]) {
        if err := cmd.RunScriptFile(os.Args[1]); err != nil {
            os.Exit(1)
        }
        return
    }
    cmd.Execute()
}
```

2. 在 `go.mod` 中添加扩展依赖：

```go
require (
    github.com/php-any/origami-openai v0.1.0
)

// 本地开发时使用 replace 指向本地路径
replace github.com/php-any/origami-openai => ../extensions/openai
```

3. 编译运行：

```bash
go build -o myapp .
./myapp script.php
```

### 内置测试

扩展自带了完整的测试套件和独立测试运行器：

```bash
cd extensions/openai
go build -o testrunner ./cmd/
./testrunner examples/run_all.php       # 运行全部测试
./testrunner examples/test_chat.php     # 单功能测试
```

### 不使用 OpenAI 扩展

如果你不需要 OpenAI 功能，只需不在 `SetRuntimeLoader` 中调用 `openai.Load(vm)` 即可。扩展模块不会编译进你的二进制文件。

## 最佳实践

1. **API Key 安全**：不要在脚本中硬编码 API Key，使用环境变量 `OPENAI_API_KEY`
2. **错误处理**：始终用 try/catch 包裹 API 调用
3. **Token 控制**：使用 `maxTokens` 限制输出长度，避免意外消耗
4. **模型选择**：根据需求选择合适的模型，如简单任务可用 `gpt-4o-mini`

## 注意事项

- OpenAI 模块需要网络访问才能正常工作
- API 调用会产生费用，请参考 [OpenAI 定价](https://openai.com/pricing)
- 语音合成的输出文件需要有写入权限
- 语音转文字支持常见音频格式（mp3、wav、m4a 等）

## 相关文档

- [Go 集成](go-integration.md) - 了解如何将 Go 库集成到折言
- [扩展开发](extensions.md) - 开发自定义扩展
- [标准库参考](./std/README.md) - 内置标准库使用说明
