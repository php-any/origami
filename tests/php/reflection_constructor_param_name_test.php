<?php

namespace tests\php;

/**
 * ReflectionClass::getConstructor()->getParameters()[n]->getName()
 * 必须返回不含 $ 的参数名。Laravel 容器 contextual 绑定用 '$'.$parameter->getName()。
 */

class ReflectionCtorParamName_Repo
{
    public $connection;

    public function __construct(string $connection, ?int $chunkSize = null)
    {
        $this->connection = $connection;
    }
}

$ref = new \ReflectionClass(ReflectionCtorParamName_Repo::class);
$ctor = $ref->getConstructor();
if ($ctor === null) {
    Log::fatal('getConstructor 不应为 null');
}

$params = $ctor->getParameters();
if (count($params) !== 2) {
    Log::fatal('构造函数参数数量错误: '.count($params));
}

$connection = $params[0];
if ($connection->getName() !== 'connection') {
    Log::fatal('getName 应为 connection, 实际 '.var_export($connection->getName(), true));
}
if (($connection->name ?? null) !== 'connection') {
    Log::fatal('public $name 应为 connection, 实际 '.var_export($connection->name ?? null, true));
}

$chunk = $params[1];
if ($chunk->getName() !== 'chunkSize') {
    Log::fatal('第二参数 getName 应为 chunkSize, 实际 '.var_export($chunk->getName(), true));
}
if (!$chunk->isDefaultValueAvailable()) {
    Log::fatal('chunkSize 应有默认值');
}
if (!$chunk->allowsNull()) {
    Log::fatal('?int $chunkSize = null 的 allowsNull 应为 true');
}

Log::info('reflection_constructor_param_name 测试通过');
