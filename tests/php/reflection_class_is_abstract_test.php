<?php

namespace tests\php;

/**
 * ReflectionClass::isAbstract / isInterface 应对齐 PHP。
 */
abstract class ReflectionIsAbstract_Abs {}
class ReflectionIsAbstract_Conc {}
interface ReflectionIsAbstract_Iface {}

$abs = (new \ReflectionClass(ReflectionIsAbstract_Abs::class))->isAbstract();
$conc = (new \ReflectionClass(ReflectionIsAbstract_Conc::class))->isAbstract();
$iface = (new \ReflectionClass(ReflectionIsAbstract_Iface::class))->isInterface();

if (!$abs) {
    Log::fatal('abstract class 应为 isAbstract=true');
}
if ($conc) {
    Log::fatal('concrete class 应为 isAbstract=false');
}
if (!$iface) {
    Log::fatal('interface 应为 isInterface=true');
}

Log::info('reflection_class_is_abstract 测试通过');
