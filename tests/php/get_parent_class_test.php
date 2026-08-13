<?php

namespace tests\php;

/**
 * get_parent_class：对象与字符串类名均可，字符串路径应能 autoload。
 */

class GetParentClass_Base {}

class GetParentClass_Child extends GetParentClass_Base {}

$obj = new GetParentClass_Child();
$parentFromObj = get_parent_class($obj);
if ($parentFromObj !== GetParentClass_Base::class && $parentFromObj !== 'tests\\php\\GetParentClass_Base') {
    Log::fatal('get_parent_class(object) 失败: ' . var_export($parentFromObj, true));
}

$parentFromName = get_parent_class(GetParentClass_Child::class);
if ($parentFromName !== GetParentClass_Base::class && $parentFromName !== 'tests\\php\\GetParentClass_Base') {
    Log::fatal('get_parent_class(string) 失败: ' . var_export($parentFromName, true));
}

$none = get_parent_class(GetParentClass_Base::class);
if ($none !== false) {
    Log::fatal('无父类应返回 false: ' . var_export($none, true));
}

Log::info('get_parent_class_test 测试通过');
