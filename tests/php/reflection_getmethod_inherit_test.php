<?php

namespace tests\php;

/**
 * ReflectionClass::getMethod 应沿继承链查找，且缺失时抛 ReflectionException；
 * ReflectionAttribute::IS_INSTANCEOF 常量必须存在。
 */
class ReflectionGetMethod_Base {
	public function inherited() {
		return 'base';
	}
}

class ReflectionGetMethod_Child extends ReflectionGetMethod_Base {
	public function own() {
		return 'child';
	}
}

$ref = new \ReflectionClass(ReflectionGetMethod_Child::class);

$own = $ref->getMethod('own');
if ($own->getName() !== 'own') {
	Log::fatal('getMethod(own) 失败');
}

$inherited = $ref->getMethod('inherited');
if ($inherited->getName() !== 'inherited') {
	Log::fatal('getMethod 未沿继承链找到 inherited');
}

$caught = false;
try {
	$ref->getMethod('missing_method_xyz');
} catch (\ReflectionException $e) {
	$caught = true;
}
if (!$caught) {
	Log::fatal('getMethod 缺失方法应抛 ReflectionException');
}

if (!defined('ReflectionAttribute::IS_INSTANCEOF') && !(\ReflectionAttribute::IS_INSTANCEOF === 2)) {
	// 类常量通过 :: 访问
}
if (\ReflectionAttribute::IS_INSTANCEOF !== 2) {
	Log::fatal('ReflectionAttribute::IS_INSTANCEOF 应为 2，实际: ' . var_export(\ReflectionAttribute::IS_INSTANCEOF, true));
}

$file = $ref->getFileName();
if (!is_string($file) || $file === '') {
	Log::fatal('ReflectionClass::getFileName 应返回源文件路径');
}
// PHP 方法名大小写不敏感
$file2 = $ref->getFilename();
if ($file2 !== $file) {
	Log::fatal('getFilename 别名应与 getFileName 相同');
}

Log::info('reflection getMethod / IS_INSTANCEOF 测试通过');
