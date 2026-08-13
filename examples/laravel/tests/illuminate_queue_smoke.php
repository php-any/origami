<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Events\Dispatcher;
use Illuminate\Foundation\Application;
use Illuminate\Queue\CallQueuedClosure;
use Illuminate\Queue\Jobs\JobName;
use Illuminate\Queue\SyncQueue;
use Illuminate\Support\Str;

/**
 * SyncQueue::push（字符串 job）+ Str::uuid + Foundation CallQueuedClosure 可加载。
 * 完整 Closure 序列化仍缺 ReflectionFunction::getFileName/getStartLine 等（另开上游）。
 */

$uuid = (string) Str::uuid();
if (!preg_match('/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i', $uuid)) {
    echo "FAIL: Str::uuid format\n";
    var_export($uuid);
    echo "\n";
    exit(1);
}

$container = new Application(dirname(__DIR__));
$container->instance('events', new Dispatcher($container));

$queue = new SyncQueue();
$queue->setContainer($container);

if ($queue->getContainer() !== $container) {
    echo "FAIL: container binding\n";
    exit(1);
}

$parsed = JobName::parse('App\\Jobs\\Demo@handle');
if (($parsed[0] ?? null) !== 'App\\Jobs\\Demo' || ($parsed[1] ?? null) !== 'handle') {
    echo "FAIL: JobName::parse\n";
    var_export($parsed);
    echo "\n";
    exit(1);
}

$ran = false;

class QueueSmoke_Job
{
    public function handle($job, $data = null)
    {
        global $ran;
        $ran = is_array($data) && ($data['x'] ?? null) === 1;
    }
}

$id = $queue->push(QueueSmoke_Job::class . '@handle', ['x' => 1]);
if ($id !== 0 && $id !== '0') {
    echo "FAIL: unexpected push id\n";
    var_export($id);
    echo "\n";
    exit(1);
}

if (!$ran) {
    echo "FAIL: SyncQueue::push job not executed\n";
    exit(1);
}

if (!class_exists(CallQueuedClosure::class)) {
    echo "FAIL: CallQueuedClosure missing\n";
    exit(1);
}

if (!class_exists(\Illuminate\Foundation\Bus\PendingDispatch::class)) {
    echo "FAIL: PendingDispatch missing\n";
    exit(1);
}

echo "PASS\n";
