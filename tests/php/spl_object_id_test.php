<?php

namespace tests\php;

/**
 * 测试 spl_object_id：同一对象返回相同整数 ID，不同对象返回不同 ID。
 */

$a = new \stdClass();
$b = new \stdClass();

$idA1 = spl_object_id($a);
$idA2 = spl_object_id($a);
$idB = spl_object_id($b);

if ($idA1 !== $idA2) {
    Log::fatal('spl_object_id_test: 同一对象应返回相同 ID');
}

if ($idA1 === $idB) {
    Log::fatal('spl_object_id_test: 不同对象应返回不同 ID');
}

if (!is_int($idA1)) {
    Log::fatal('spl_object_id_test: 返回值应为 int');
}

Log::info('spl_object_id 测试通过');
