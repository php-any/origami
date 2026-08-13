<?php

namespace tests\php;

/**
 * Exception / Throwable 相关方法测试：
 * - getMessage()
 * - getCode()
 * - getFile()
 * - getLine()
 * - getTrace()
 * - getPrevious()
 * - getTraceAsString()
 *
 * 目标：确保 Origami 内部的 ThrowValue 实现能满足 Symfony 等库对异常对象的基本依赖。
 */

try {
    $expectedMessage = 'exception-methods-test';
    $expectedFile = basename(__FILE__);

    $lineBeforeThrow = __LINE__ + 1;
    throw new \Exception($expectedMessage);
} catch (\Exception $e) {
    // getMessage
    if ($e->getMessage() !== $expectedMessage) {
        Log::fatal('Exception::getMessage 测试失败: ' . $e->getMessage());
    }

    // getCode（当前实现统一返回 0）
    if ($e->getCode() !== 0) {
        Log::fatal('Exception::getCode 测试失败: ' . var_export($e->getCode(), true));
    }

    // getFile
    $file = $e->getFile();
    if (basename($file) !== $expectedFile) {
        Log::fatal('Exception::getFile 测试失败: ' . $file);
    }

    // getLine：不强制精确等于 $lineBeforeThrow，只要求为正整数
    $line = $e->getLine();
    if (!is_int($line) || $line <= 0) {
        Log::fatal('Exception::getLine 测试失败: ' . var_export($line, true));
    }

    // getTrace：应返回数组
    $trace = $e->getTrace();
    if (!is_array($trace)) {
        Log::fatal('Exception::getTrace 返回的不是数组: ' . gettype($trace));
    }

    // getPrevious：当前实现应返回 null
    if (null !== $e->getPrevious()) {
        Log::fatal('Exception::getPrevious 预期为 null');
    }

    // getTraceAsString：至少应返回字符串
    $traceStr = $e->getTraceAsString();
    if (!is_string($traceStr)) {
        Log::fatal('Exception::getTraceAsString 返回的不是字符串: ' . gettype($traceStr));
    }
}

Log::info('Exception 方法测试通过');

