<?php

namespace tests\php;

/**
 * preg_match 命名捕获组：同时提供命名键与数值键；array_slice 后命名键仍可用。
 */

$subject = '/telescope/telescope-api/requests/754132c7-2e36-41df-bb75-228ba11567d5';
$pattern = '{^/telescope/telescope\-api/requests/(?P<telescopeEntryId>[^/]++)$}sDu';

$matches = [];
$ok = preg_match($pattern, $subject, $matches);
if ($ok !== 1) {
    Log::fatal('preg_match_named_groups: 期望匹配成功');
}

if (!isset($matches['telescopeEntryId']) || $matches['telescopeEntryId'] !== '754132c7-2e36-41df-bb75-228ba11567d5') {
    Log::fatal('preg_match_named_groups: 缺少命名键 telescopeEntryId');
}

if (!isset($matches[1]) || $matches[1] !== '754132c7-2e36-41df-bb75-228ba11567d5') {
    Log::fatal('preg_match_named_groups: 缺少数值键 1');
}

$keys = array_keys($matches);
if ($keys[0] !== 0 || $keys[1] !== 'telescopeEntryId' || $keys[2] !== 1) {
    Log::fatal('preg_match_named_groups: array_keys 顺序不符合 PHP: ' . json_encode($keys));
}

$sliced = array_slice($matches, 1);
$parameters = array_intersect_key($sliced, array_flip(['telescopeEntryId']));
if (($parameters['telescopeEntryId'] ?? null) !== '754132c7-2e36-41df-bb75-228ba11567d5') {
    Log::fatal('preg_match_named_groups: array_slice + array_intersect_key 丢失命名键');
}

Log::info('preg_match_named_groups 测试通过');
