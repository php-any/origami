<?php

namespace tests\php;

/**
 * ErrorException 6 参构造应对齐 PHP：previous 是 Throwable，不是 severity int。
 * Laravel CompilerEngine：new ViewException($msg, 0, 1, $file, $line, $e)
 * 若 previous 被写成 int，ViewException::report → Reflector::isCallable
 * 会 new ReflectionMethod($class, '__callStatic') 并 TypeError。
 */

class ErrorExCtor_ViewException extends \ErrorException
{
}

$inner = new \Exception('blade failed');
$file = __FILE__;
$line = 42;

$e = new \ErrorException($inner->getMessage() . ' (View: x)', 0, 1, $file, $line, $inner);

$prev = $e->getPrevious();
if (!($prev instanceof \Exception)) {
    \Log::fatal('ErrorException::getPrevious 应为原始异常, 实际: ' . gettype($prev) . ' ' . var_export($prev, true));
}
if ($prev->getMessage() !== 'blade failed') {
    \Log::fatal('getPrevious 消息不对: ' . $prev->getMessage());
}
if ($e->getSeverity() !== 1) {
    \Log::fatal('getSeverity 应为 1, 实际: ' . var_export($e->getSeverity(), true));
}
if ($e->getFile() !== $file) {
    \Log::fatal('getFile 应为传入路径, 实际: ' . $e->getFile());
}
if ($e->getLine() !== $line) {
    \Log::fatal('getLine 应为 42, 实际: ' . var_export($e->getLine(), true));
}

$wrapped = new ErrorExCtor_ViewException($inner->getMessage() . ' (View: x)', 0, 1, $file, $line, $inner);
$wprev = $wrapped->getPrevious();
if (!($wprev instanceof \Exception) || $wprev->getMessage() !== 'blade failed') {
    \Log::fatal('子类 ViewException::getPrevious 应为原始异常, 实际: ' . gettype($wprev) . ' ' . var_export($wprev, true));
}

// 对齐 Illuminate\Support\Reflector::isCallable([$exception, 'report'])
$var = [$wprev, 'report'];
if (!is_object($var[0]) || !is_string($var[1])) {
    \Log::fatal('callable 数组 [previous, report] 形态不对');
}
$class = is_object($var[0]) ? get_class($var[0]) : $var[0];
if (!is_string($class)) {
    \Log::fatal('Reflector $class 应为 string, 实际 ' . gettype($class));
}
if (!class_exists($class)) {
    \Log::fatal('class_exists($class) 应为 true: ' . $class);
}

if (method_exists($class, 'report')) {
    $rm = new \ReflectionMethod($class, 'report');
    if ($rm->getName() !== 'report') {
        \Log::fatal('ReflectionMethod(report) 失败');
    }
} elseif (!is_object($var[0]) && method_exists($class, '__callStatic')) {
    \Log::fatal('previous 不是对象却走到 __callStatic，会把非 string 传给 ReflectionMethod');
} else {
    try {
        $rm = new \ReflectionMethod($class, '__callStatic');
        \Log::fatal('Exception 不应有公开 __callStatic: ' . $rm->getName());
    } catch (\Throwable $t) {
        // PHP：方法不存在时 ReflectionMethod 抛 ReflectionException，不是 TypeError(int)
        if (strpos($t->getMessage(), 'IntValue') !== false || strpos($t->getMessage(), 'int given') !== false) {
            \Log::fatal('ReflectionMethod 仍收到 int: ' . $t->getMessage());
        }
    }
}

\Log::info('ErrorException 6 参构造 / getPrevious 测试通过');
