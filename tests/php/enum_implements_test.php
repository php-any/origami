<?php

namespace tests\php;

/**
 * enum Name: string implements Interface 语法（Filament Heroicon 依赖）。
 */
interface EnumImplements_Iface
{
    public function label(): string;
}

enum EnumImplements_Status: string implements EnumImplements_Iface
{
    case Open = 'open';
    case Closed = 'closed';

    public function label(): string
    {
        return $this->value;
    }
}

if (!enum_exists(EnumImplements_Status::class) && !class_exists(EnumImplements_Status::class)) {
    Log::fatal('enum 类未加载');
}

$v = EnumImplements_Status::Open;
if ($v->value !== 'open') {
    Log::fatal('case 值错误');
}
if ($v->label() !== 'open') {
    Log::fatal('implements 方法调用失败');
}

Log::info('enum_implements 测试通过');
