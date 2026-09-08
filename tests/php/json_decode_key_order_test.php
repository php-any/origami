<?php

namespace tests\php;

/**
 * json_decode 必须保留对象键插入顺序，否则 Livewire checksum（hash_hmac(json_encode)）会失败。
 */

$json = '{"data":{"email":"","password":"","remember":false},"memo":{"id":"x","name":"admin.login"},"checksum":"abc"}';
$arr = json_decode($json, associative: true);
$keys = array_keys($arr);
if ($keys !== ['data', 'memo', 'checksum']) {
    Log::fatal('顶层键序错误: ' . json_encode($keys));
}
$dataKeys = array_keys($arr['data']);
if ($dataKeys !== ['email', 'password', 'remember']) {
    Log::fatal('data 键序错误: ' . json_encode($dataKeys));
}
$re = json_encode($arr);
if ($re !== $json) {
    Log::fatal("roundtrip 键序/内容不一致\norig=$json\nre  =$re");
}

Log::info('json_decode 键序保留测试通过');
