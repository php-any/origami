<?php

namespace tests\php;

/**
 * ReflectionProperty::IS_* 是类常量，:: 访问必须得到 PHP 官方整数值。
 */

if (ReflectionProperty::IS_PUBLIC !== 1) {
    \Log::fatal('IS_PUBLIC 应为 1，实际 ' . var_export(ReflectionProperty::IS_PUBLIC, true));
}
if (ReflectionProperty::IS_PROTECTED !== 2) {
    \Log::fatal('IS_PROTECTED 应为 2，实际 ' . var_export(ReflectionProperty::IS_PROTECTED, true));
}
if (ReflectionProperty::IS_PRIVATE !== 4) {
    \Log::fatal('IS_PRIVATE 应为 4，实际 ' . var_export(ReflectionProperty::IS_PRIVATE, true));
}
if (ReflectionProperty::IS_STATIC !== 16) {
    \Log::fatal('IS_STATIC 应为 16，实际 ' . var_export(ReflectionProperty::IS_STATIC, true));
}
if ((ReflectionProperty::IS_PUBLIC | ReflectionProperty::IS_PROTECTED | ReflectionProperty::IS_PRIVATE) !== 7) {
    \Log::fatal('可见性位或结果错误');
}

\Log::info('reflection_property_constants 测试通过');
