<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Broadcasting\BroadcastManager;
use Illuminate\Config\Repository;
use Illuminate\Container\Container;
use Psr\Log\LoggerInterface;

/**
 * 使用 log 驱动，不依赖 Pusher / Redis / Ably。
 */
class BroadcastingSmoke_StubLogger implements LoggerInterface
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
        $this->entries[] = [(string) $level, (string) $message];
    }
}

$app = new Container();
$app->instance('config', new Repository([
    'broadcasting' => [
        'default' => 'log',
        'connections' => [
            'log' => ['driver' => 'log'],
        ],
    ],
]));
$logger = new BroadcastingSmoke_StubLogger();
$app->instance(LoggerInterface::class, $logger);

$manager = new BroadcastManager($app);
$broadcaster = $manager->connection('log');
$broadcaster->broadcast(['orders'], 'OrderPlaced', ['id' => 7]);

if (count($logger->entries) !== 1) {
    echo "FAIL: log entry count\n";
    exit(1);
}

[$level, $message] = $logger->entries[0];
if ($level !== 'info' || strpos($message, 'OrderPlaced') === false || strpos($message, 'orders') === false) {
    echo "FAIL: broadcast log\n";
    var_export($logger->entries);
    echo "\n";
    exit(1);
}

echo "PASS\n";
