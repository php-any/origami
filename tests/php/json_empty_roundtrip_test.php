<?php

namespace tests\php;

/**
 * json_encode 空数组应为 []；空对象/空关联结构往返需与 checksum 语义一致。
 */
$a = json_encode([]);
if ($a !== '[]') {
    Log::fatal('空列表 encode 失败: '.$a);
}

$b = json_decode('{}', associative: true);
$be = json_encode($b);
// PHP: json_decode('{}', true) 再 encode 为 []
if ($be !== '[]') {
    Log::fatal("空对象 decode associative 再 encode 应为 [], got=$be");
}

$c = json_decode('{"errors":{}}', associative: true);
$ce = json_encode($c);
if ($ce !== '{"errors":[]}') {
    Log::fatal("errors 空对象往返: $ce");
}

Log::info('json_empty_roundtrip 测试通过: a='.$a.' be='.$be.' ce='.$ce);
