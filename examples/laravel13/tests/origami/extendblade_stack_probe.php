<?php
/**
 * 最小探针：ExtendBlade 栈 + Notifications mount，不走完整 HTTP。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$ref = new ReflectionClass(Livewire\Mechanisms\ExtendBlade\ExtendBlade::class);
$prop = $ref->getProperty('livewireComponents');
$prop->setAccessible(true);

echo "stack_initial=".json_encode($prop->getValue())."\n";

// 手动 push/pop 验证 static 属性
$dummy = new stdClass();
$dummy->x = 1;
// 通过公开 API
$ebClass = Livewire\Mechanisms\ExtendBlade\ExtendBlade::class;
// startLivewireRendering 是实例方法；取 mechanism 实例
$mech = null;
try {
    $manager = app('livewire');
    echo "livewire_ok\n";
} catch (Throwable $e) {
    echo "livewire_fail: ".$e->getMessage()."\n";
}

Livewire\before('render', function ($target, $view) use ($prop) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    $stack = $prop->getValue();
    $n = is_array($stack) ? count($stack) : -1;
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "before render target={$cls} current={$curCls} stackN={$n}\n",
        FILE_APPEND
    );
});

Livewire\after('render', function ($target, $view) use ($prop) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "after render target={$cls} current={$curCls}\n",
        FILE_APPEND
    );
});

@unlink(storage_path('framework/stack-probe.txt'));

echo "mount_notifications...\n";
try {
    $html = app('livewire')->mount(\Filament\Notifications\Livewire\Notifications::class);
    echo "notif_len=".strlen((string)$html)."\n";
    echo "has_err=".(str_contains((string)$html, 'getBroadcastChannel') ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "---stack-probe---\n";
echo @file_get_contents(storage_path('framework/stack-probe.txt'));
echo "stack_final=".json_encode($prop->getValue())."\n";
echo "DONE\n";
