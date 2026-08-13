<?php

namespace tests\php;

/**
 * 接口多继承测试：
 * - MultiExt_Base1
 * - MultiExt_Base2
 * - MultiExt_Sub extends MultiExt_Base1, MultiExt_Base2
 * - MultiExt_Impl implements MultiExt_Sub
 *
 * 验证：
 * 1) instanceof 对多继承所有父接口生效
 * 2) 接口类型提示支持多继承（data.Class.Is / type_class.go + parser/interface_parser.go）
 */

interface MultiExt_Base1
{
}

interface MultiExt_Base2
{
}

interface MultiExt_Sub extends MultiExt_Base1, MultiExt_Base2
{
}

class MultiExt_Impl implements MultiExt_Sub
{
}

$impl = new MultiExt_Impl();

// 1. instanceof + 接口多继承
if (!($impl instanceof MultiExt_Sub)) {
    Log::fatal("接口多继承 instanceof(Sub) 测试失败：期望 true");
}
Log::info("接口多继承 instanceof(Sub) 测试通过");

if (!($impl instanceof MultiExt_Base1)) {
    Log::fatal("接口多继承 instanceof(Base1) 测试失败：期望 true");
}
Log::info("接口多继承 instanceof(Base1) 测试通过");

if (!($impl instanceof MultiExt_Base2)) {
    Log::fatal("接口多继承 instanceof(Base2) 测试失败：期望 true");
}
Log::info("接口多继承 instanceof(Base2) 测试通过");

// 2. 接口类型提示 + 运行时类型检查（依赖 data.Class.Is / type_class.go）
function MultiExt_acceptBase1(MultiExt_Base1 $x): void
{
    Log::info("接口多继承 类型提示 Base1 测试通过: " . get_class($x));
}

function MultiExt_acceptBase2(MultiExt_Base2 $x): void
{
    Log::info("接口多继承 类型提示 Base2 测试通过: " . get_class($x));
}

MultiExt_acceptBase1($impl);
MultiExt_acceptBase2($impl);

