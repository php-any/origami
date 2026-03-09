<?php

namespace tests\php;

/**
 * class_alias 函数测试（当前为最小实现）：
 * - 对已存在类调用 class_alias 返回 true
 * - 对不存在类调用 class_alias 返回 false
 *
 * 注意：当前 Go 端实现暂未真正把别名绑定到原始类，只验证基本返回语义。
 */

class ClassAlias_Original {}

// 1. 已存在类，期望返回 true
if (!class_alias(ClassAlias_Original::class, 'ClassAlias_Alias')) {
    Log::fatal('class_alias(ClassAlias_Original::class, ClassAlias_Alias) 期望返回 true');
}
Log::info('class_alias 已存在类测试通过');

// 2. 不存在类，期望返回 false
if (class_alias('ClassAlias_NotDefined', 'ClassAlias_Other')) {
    Log::fatal('class_alias(ClassAlias_NotDefined, ClassAlias_Other) 期望返回 false');
}
Log::info('class_alias 不存在类测试通过');

