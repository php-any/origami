<?php

namespace tests\php;

/**
 * 验证方法返回类型交集语法 A & B 可被解析，且 Reflection 可见。
 */
interface IntersectionReturnType_A {}
interface IntersectionReturnType_B {}

class IntersectionReturnType_Both implements IntersectionReturnType_A, IntersectionReturnType_B {}

class IntersectionReturnType_Holder
{
    public function getBoth(): IntersectionReturnType_A & IntersectionReturnType_B
    {
        return new IntersectionReturnType_Both();
    }
}

$m = new \ReflectionMethod(IntersectionReturnType_Holder::class, 'getBoth');
$rt = $m->getReturnType();
if ($rt === null) {
    Log::fatal('交集返回类型未被解析');
}

$obj = new IntersectionReturnType_Holder();
$v = $obj->getBoth();
if (!($v instanceof IntersectionReturnType_A) || !($v instanceof IntersectionReturnType_B)) {
    Log::fatal('交集返回值实例检查失败');
}

Log::info('intersection_return_type 测试通过');
