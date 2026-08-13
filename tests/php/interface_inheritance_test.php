<?php

namespace tests\php;

/**
 * 接口继承测试：
 * - InterfaceInheritance_Base
 * - InterfaceInheritance_Mid extends InterfaceInheritance_Base
 * - InterfaceInheritance_Sub extends InterfaceInheritance_Mid
 * - InterfaceInheritance_Impl implements InterfaceInheritance_Sub
 *
 * 验证：
 * 1) instanceof 对接口继承链生效（Impl instanceof Base/Mid/Sub 都应为 true）
 * 2) 接口类型提示的运行时类型检查支持接口继承：
 *    - function InterfaceInheritance_acceptBase(InterfaceInheritance_Base $x)
 *      允许传入实现了子接口的实例。
 */

interface InterfaceInheritance_Base {}

interface InterfaceInheritance_Mid extends InterfaceInheritance_Base {}

interface InterfaceInheritance_Sub extends InterfaceInheritance_Mid {}

class InterfaceInheritance_Impl implements InterfaceInheritance_Sub {}

$impl = new InterfaceInheritance_Impl();

// 1. instanceof + 接口继承链
if ($impl instanceof InterfaceInheritance_Sub) {
    Log::info("接口继承 instanceof(Sub) 测试通过");
} else {
    Log::fatal("接口继承 instanceof(Sub) 测试失败：期望 true");
}

if ($impl instanceof InterfaceInheritance_Mid) {
    Log::info("接口继承 instanceof(Mid) 测试通过");
} else {
    Log::fatal("接口继承 instanceof(Mid) 测试失败：期望 true");
}

if ($impl instanceof InterfaceInheritance_Base) {
    Log::info("接口继承 instanceof(Base) 测试通过");
} else {
    Log::fatal("接口继承 instanceof(Base) 测试失败：期望 true");
}

// 2. 接口类型提示 + 运行时类型检查（依赖 data.Class.Is / type_class.go）
function InterfaceInheritance_acceptBase(InterfaceInheritance_Base $x): void
{
    Log::info("接口继承 类型提示测试通过: " . get_class($x));
}

InterfaceInheritance_acceptBase($impl);

