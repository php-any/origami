<?php

namespace tests\php;

/**
 * Exception::$message 等受保护属性：子类可直接读写（Laravel ModelNotFoundException::setModel）
 */

class ExceptionMessageProp_Custom extends \RuntimeException
{
    public function setMsg($model, $ids = [])
    {
        $this->message = "No query results for model [{$model}]";
        if (count($ids) > 0) {
            $this->message .= ' ' . implode(', ', $ids);
        } else {
            $this->message .= '.';
        }
        return $this;
    }
}

try {
    throw (new ExceptionMessageProp_Custom())->setMsg('App\\User', [1, 2]);
} catch (\Throwable $e) {
    $msg = $e->getMessage();
    $expect = 'No query results for model [App\\User] 1, 2';
    if ($msg !== $expect) {
        Log::fatal('getMessage 与 $this->message 不同步: got=' . $msg);
    }
}

$e2 = (new ExceptionMessageProp_Custom())->setMsg('X');
if ($e2->getMessage() !== 'No query results for model [X].') {
    Log::fatal('空 ids 时 message .= "." 失败: ' . $e2->getMessage());
}

$base = new \Exception('hello');
if ($base->getMessage() !== 'hello') {
    Log::fatal('普通 Exception::getMessage 回归失败');
}

Log::info('Exception message 属性测试通过');
