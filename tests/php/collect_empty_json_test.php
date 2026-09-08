<?php

namespace tests\php;

$payload = [
    'errors' => collect([])
        ->filter(function ($value, $key) {
            return false;
        })
        ->toArray(),
];
$enc = json_encode($payload);
if ($enc !== '{"errors":[]}') {
    Log::fatal('collect filter toArray empty encode 失败: '.$enc);
}

$payload2 = ['errors' => []];
$enc2 = json_encode($payload2);
if ($enc2 !== '{"errors":[]}') {
    Log::fatal('plain empty encode 失败: '.$enc2);
}

Log::info('collect_empty_json 测试通过');
