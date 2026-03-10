---
name: param-types-docs
description: Document function parameter type syntax and data.Types constraints in the Origami runtime. Use when the user asks about parameter types, union/nullable types, or data package type rules, or wants to update related documentation.
---

# 函数参数类型与 `data.Types` 文档

## 使用时机

当用户：
- 提到「函数参数类型」「联合类型 string|int|null」「?nullable 类型」「类型约束」「`data.Types`」等关键词时
- 想梳理/补充关于**函数参数类型语法**与 **`data` 包类型系统**（`Types`、`UnionType`、`NullableType` 等）的文档
- 在阅读/修改 `parser/parameter_parser.go`、`data/types.go` 等文件并希望同步文档

就应用本 Skill 的指引来回答或编辑文档。

## 核心概念速览

- **参数类型语法（解析侧）**
  - 由 `parser/parameter_parser.go` 中的 `parseSingleParameter` 和 `parseConstructorParameterType` 负责解析：
    - 支持**联合类型**：`string|int|null`
    - 支持**可空类型前缀**：`?string`
    - 支持 `null`、`false` 作为类型参与联合
    - 支持「类型在前」和「变量在前」两种风格，例如：
      - `string $data`、`?User $u`、`string|int $v`
      - `data: string`（`name: type` 形式）
    - 支持可变参数与引用参数的组合：`type ...$vars`、`&...$vars`、`type &...$vars`
  - 解析完成后，会为参数构造一个 `data.Types` 实例，并存入当前作用域，用于运行时检查与反射（`ReflectionParameter::getType()` 等）。

- **类型系统与约束（运行时侧）**
  - `data.Types` 接口定义在 `data/types.go`：
    - `Is(value Value) bool`：判断某个 `Value` 是否满足该类型
    - `String() string`：类型标识字符串（不带泛型信息）
  - 重要实现：
    - `BaseType` 系列通过工厂函数 `NewBaseType(ty string)` 创建：
      - 基础类型：`int`、`string`、`bool`、`float`、`array`、`object`、`callable` 等
      - 特殊类型：`false`（等同 `bool`）、`null`（`NullType`）、`self`/`static`（`StaticType`）、`closure`/`\Closure`（`ClosureType`）
      - 联合与可空在字符串层先做拆分处理：
        - 包含 `|` 时构造 `UnionType`
        - 以 `?` 开头时构造 `NullableType`
      - 其它字符串按类名处理，生成 `Class{Name: ty}`
    - `UnionType`：表示 `type1|type2|...`，`Is` 对内部所有类型做「任意一个为真即可」判断
    - `NullableType`：包装一个基础类型，允许值为 `null` 或基础类型值
    - `MultipleReturnType`：用于多返回值检测（数组长度与每个元素逐一匹配）
    - `NewGenericType`：用于泛型类型（如 `Generic{Name, Types}`），在文档中可单独说明。
  - `ISBaseType(ty string)` 用来判断某个字符串是否是「内建基础类型关键字」，影响 `parserType` 如何解析类名 vs 基础类型。

## 回答设计（参数类型相关问题）

当用户问「参数类型怎么写/怎么解析/怎么约束」时，建议用下面的结构：

1. **先解释语法层**（指向 `parser/parameter_parser.go`）
   - 支持哪些写法（联合、可空、`name: type`、引用、可变参数等）
   - 这些写法最终如何映射到 `data.Types`
2. **再解释运行时类型系统**（指向 `data/types.go`）
   - 说明 `Types` 接口、`NewBaseType`、`UnionType`、`NullableType` 等的行为
   - 给出常见基础类型映射表
3. **最后给出限制与建议**
   - 哪些字符串会被当作类名
   - `void`/`mixed` 的特殊处理（`NewBaseType` 返回 `nil`）
   - 对多返回值、泛型的约束边界

### 示例回答模版

用户问：「函数参数的联合类型是怎么实现的？」

可以按照如下模版回答（根据上下文裁剪）：

```markdown
**语法解析**：在 `parser/parameter_parser.go` 里，`parseSingleParameter` 会识别形如 `string|int|null $v` 的写法，内部通过 `parseConstructorParameterType` 依次解析每个片段，构造一个 `[]data.Types`。

**类型系统**：这些 `data.Types` 片段会被打包成一个 `data.UnionType`（由 `data.NewUnionType` 创建），其 `Is(value)` 方法会对内部每个类型逐一调用 `Is`，任意一个满足即可。

**可空/特殊值**：`null` 和 `false` 也可以出现在联合类型中，对应 `NullType` 和等同于 `Bool` 的类型，实现上通过 `NewBaseType("null")` 与 `NewBaseType("false")` 完成。
```

## 文档编写/补充建议

当在 `docs/` 或其他地方补充参数类型与 `data.Types` 文档时，推荐结构：

1. **语法说明章节**
   - 对应 Parser 侧：列出支持的参数声明形式及例子
   - 特别说明联合类型、可空类型、`name: type`、引用与可变参数的组合
2. **类型系统章节**
   - 对应 `data` 侧：解释 `Types` 接口，实现类（`UnionType`、`NullableType`、`MultipleReturnType`、`StaticType`、`ClosureType`、`NullType` 等）
   - 给一个「PHP 风格类型 → `data.Types` 实现」的映射表
3. **约束与边界**
   - 哪些字符串会被当作类名，哪些是基础类型
   - `void`、`mixed` 返回 `nil` 的语义（只用于返回值、不做运行时检查）
   - 多返回值/泛型目前支持到什么程度

后续在整理具体文档时，可以直接复用这里的结构与术语，保持项目中关于参数类型与 `data` 类型系统的说明风格一致。
