<?php

namespace tests\php;

/**
 * enum 内 public const 应可通过 EnumName::CONST 访问。
 */

enum EnumConst_Level: int
{
    case Debug = 100;
    case Info = 200;

    public const NAMES = ['DEBUG', 'INFO'];
    public const VALUES = [100, 200];
}

$names = EnumConst_Level::NAMES;
if (!is_array($names) || $names[0] !== 'DEBUG') {
    Log::fatal('EnumConst_Level::NAMES 错误: ' . var_export($names, true));
}

$vals = EnumConst_Level::VALUES;
if (!is_array($vals) || $vals[0] !== 100) {
    Log::fatal('EnumConst_Level::VALUES 错误: ' . var_export($vals, true));
}

$merged = EnumConst_Level::NAMES + EnumConst_Level::VALUES;
if (count($merged) < 2) {
    Log::fatal('NAMES + VALUES 合并失败');
}

Log::info('enum_const_static 测试通过');
