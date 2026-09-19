<?php

namespace tests\php;

/**
 * PDOException::getTraceAsString() 必须返回 string。
 * Symfony FlattenException 会把它赋给 private string $traceAsString。
 */

try {
    throw new \PDOException('pdo-trace-as-string');
} catch (\PDOException $e) {
    $traceStr = $e->getTraceAsString();
    if (!is_string($traceStr)) {
        Log::fatal('PDOException::getTraceAsString 返回的不是字符串: ' . gettype($traceStr));
    }
    $trace = $e->getTrace();
    if (!is_array($trace)) {
        Log::fatal('PDOException::getTrace 返回的不是数组: ' . gettype($trace));
    }
    if ($e->getMessage() !== 'pdo-trace-as-string') {
        Log::fatal('PDOException::getMessage 不符合预期: ' . $e->getMessage());
    }
}

class PdoException_FlattenTraceAsString
{
    private string $traceAsString;

    public function take(\PDOException $e): void
    {
        $this->traceAsString = $e->getTraceAsString();
    }

    public function get(): string
    {
        return $this->traceAsString;
    }
}

try {
    throw new \PDOException('pdo-flatten-assign');
} catch (\PDOException $e) {
    $flatten = new PdoException_FlattenTraceAsString();
    $flatten->take($e);
    if ($flatten->get() === '') {
        Log::fatal('FlattenException 风格的 string $traceAsString 赋值为空');
    }
}

Log::info('PDOException getTraceAsString 测试通过');
