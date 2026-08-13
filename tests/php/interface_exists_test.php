<?php

namespace tests\php;

/**
 * interface_exists 函数测试：
 * - 已定义的接口返回 true
 * - 未定义的接口返回 false
 * - 结合接口继承与实现，确保自动加载路径与接口注册正常工作
 */

interface InterfaceExists_Base {}

interface InterfaceExists_Sub extends InterfaceExists_Base {}

class InterfaceExists_Impl implements InterfaceExists_Sub {}

// 1. 已存在接口（当前文件中定义）
if (!interface_exists(InterfaceExists_Base::class)) {
    Log::fatal('interface_exists(InterfaceExists_Base::class) 期望返回 true');
}
Log::info('interface_exists 已定义接口测试通过: InterfaceExists_Base');

if (!interface_exists(InterfaceExists_Sub::class)) {
    Log::fatal('interface_exists(InterfaceExists_Sub::class) 期望返回 true');
}
Log::info('interface_exists 已定义接口测试通过: InterfaceExists_Sub');

// 2. 未定义接口
if (interface_exists('InterfaceExists_NotDefined')) {
    Log::fatal('interface_exists(InterfaceExists_NotDefined) 期望返回 false');
}
Log::info('interface_exists 未定义接口测试通过');

