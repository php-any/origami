<?php

/**
 * 单步 Symfony 组件验证（子进程隔离）
 * 用法: go run zy.go examples/symfony/check_one.php step01_http_foundation.php
 */

$baseDir = __DIR__;
require $baseDir . '/vendor/autoload.php';

$stepFile = $argv[1] ?? '';
if ($stepFile === '' || !str_ends_with($stepFile, '.php')) {
    Log::fatal('用法: check_one.php <step_file.php>');
}

$pass = 0;
$fail = 0;

function step_check(string $name, bool $ok, ?string $detail = null): void
{
    global $pass, $fail;
    if ($ok) {
        $pass++;
        $msg = "[PASS] $name";
        if ($detail !== null) {
            $msg .= " — $detail";
        }
        Log::info($msg);
    } else {
        $fail++;
        $msg = "[FAIL] $name";
        if ($detail !== null) {
            $msg .= " — $detail";
        }
        Log::error($msg);
    }
}

Log::info('====== ' . $stepFile . ' ======');

try {
    include $baseDir . '/steps/' . $stepFile;
} catch (\Throwable $e) {
    step_check(basename($stepFile, '.php'), false, $e->getMessage());
    Log::info('[TRACE] ' . $e->getTraceAsString());
}

Log::info("====== $pass 通过, $fail 失败 ======");
if ($fail > 0) {
    exit(1);
}
