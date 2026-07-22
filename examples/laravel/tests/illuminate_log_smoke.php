<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Log\Logger;
use Psr\Log\LoggerInterface;

/**
 * 不依赖 Monolog 文件写入（需 DateTimeImmutable::format 等上游能力）。
 * 仅验证 illuminate/log 的 Logger 包装层与 PSR-3 委托。
 */
class LogSmoke_StubLogger implements LoggerInterface
{
    public array $entries = [];

    public function emergency($message, array $context = []): void
    {
        $this->log('emergency', $message, $context);
    }

    public function alert($message, array $context = []): void
    {
        $this->log('alert', $message, $context);
    }

    public function critical($message, array $context = []): void
    {
        $this->log('critical', $message, $context);
    }

    public function error($message, array $context = []): void
    {
        $this->log('error', $message, $context);
    }

    public function warning($message, array $context = []): void
    {
        $this->log('warning', $message, $context);
    }

    public function notice($message, array $context = []): void
    {
        $this->log('notice', $message, $context);
    }

    public function info($message, array $context = []): void
    {
        $this->log('info', $message, $context);
    }

    public function debug($message, array $context = []): void
    {
        $this->log('debug', $message, $context);
    }

    public function log($level, $message, array $context = []): void
    {
        $this->entries[] = [$level, (string) $message, $context];
    }
}

$stub = new LogSmoke_StubLogger();
$logger = new Logger($stub);
$logger->info('hello illuminate/log', ['ok' => true]);

if (count($stub->entries) !== 1) {
    echo "FAIL: expected 1 log entry\n";
    exit(1);
}

[$level, $message, $context] = $stub->entries[0];
if ($level !== 'info' || $message !== 'hello illuminate/log' || ($context['ok'] ?? null) !== true) {
    echo "FAIL: unexpected entry\n";
    var_export($stub->entries);
    echo "\n";
    exit(1);
}

echo "PASS\n";
