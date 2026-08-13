<?php

namespace tests\php;

/**
 * 接口方法返回值联合类型解析测试：
 *
 * interface InterfaceReturnTypeUnion
 * {
 *     public function parse(mixed $value): int|string;
 * }
 *
 * 目的：验证 InterfaceParser 能正确解析 `int|string` 这样的联合返回类型，
 * 不再在 `:` 之后只接受单一类型。
 */

interface InterfaceReturnTypeUnion
{
    public function parse(mixed $value): int|string;
}

// 只要脚本能正常运行到这里而不抛解析错误，就说明联合返回类型已被正确解析。
Log::info("interface return type union 解析测试通过");

