<?php

namespace tests\php;

/**
 * 子接口可访问父接口常量（CarbonInterface::SUNDAY 场景）。
 */

interface InterfaceConstParent_Base
{
    public const SUNDAY = 0;
    public const MONDAY = 1;
}

interface InterfaceConstParent_Child extends InterfaceConstParent_Base
{
}

if (InterfaceConstParent_Child::SUNDAY !== 0) {
    \Log::fatal('子接口应继承父接口常量 SUNDAY');
}
if (InterfaceConstParent_Base::MONDAY !== 1) {
    \Log::fatal('父接口常量 MONDAY 读取失败');
}

\Log::info('interface_parent_const 测试通过');
