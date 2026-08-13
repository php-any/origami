<?php

namespace tests\php;

/**
 * 多行 interface extends 列表（ramsey/uuid UuidInterface 风格）解析。
 */

interface InterfaceMultiline_BaseA
{
    public function a(): void;
}

interface InterfaceMultiline_BaseB
{
    public function b(): void;
}

interface InterfaceMultiline_MultiExtends extends
    InterfaceMultiline_BaseA,
    InterfaceMultiline_BaseB
{
    public function c(): void;
}

if (!interface_exists(InterfaceMultiline_MultiExtends::class)) {
    Log::fatal('multi-line interface extends 未注册');
}

Log::info('multi-line interface extends 测试通过');
