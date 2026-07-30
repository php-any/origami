<?php

namespace tests\php;

enum DefConst_Level: int {
    case Emergency = 600;
    case Debug = 100;
}

$name = DefConst_Level::class . '::Emergency';
var_dump($name);
var_dump(defined($name));
if (!defined($name)) {
    Log::fatal('defined(Enum::Case) 应为 true');
}
$c = constant($name);
if (!($c instanceof DefConst_Level)) {
    Log::fatal('constant(Enum::Case) 应返回 case 实例');
}
Log::info('enum_defined_constant 测试通过');
