<?php
/**
 * ExtendBlade 栈探针（不用 Reflection 读 protected static）。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

@unlink(storage_path('framework/stack-probe.txt'));

Livewire\before('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    $is = Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent() ? 'yes' : 'no';
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "BEFORE target={$cls} current={$curCls} isLW={$is}\n",
        FILE_APPEND
    );
});

// 与 ExtendBlade::on 同相位：在 listeners 里靠后，应能看到 push 后的栈
Livewire\on('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    $is = Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent() ? 'yes' : 'no';
    // 检查本监听器 $this —— 无绑定闭包应无 $this；这里用 use 外层即可
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "ON target={$cls} current={$curCls} isLW={$is}\n",
        FILE_APPEND
    );
});

Livewire\after('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "AFTER target={$cls} current={$curCls}\n",
        FILE_APPEND
    );
    return function ($html) use ($cls) {
        $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
        $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
        file_put_contents(
            storage_path('framework/stack-probe.txt'),
            "FINISH targetWas={$cls} current={$curCls}\n",
            FILE_APPEND
        );
        return $html;
    };
});

echo "mount_notifications...\n";
try {
    $html = app('livewire')->mount(\Filament\Notifications\Livewire\Notifications::class);
    echo "notif_len=".strlen((string)$html)."\n";
    echo "has_err=".(str_contains((string)$html, 'getBroadcastChannel') ? 'yes' : 'no')."\n";
    echo "this_in_html=".(str_contains((string)$html, 'Login::getBroadcastChannel') ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}

echo "---stack---\n";
echo @file_get_contents(storage_path('framework/stack-probe.txt'));
echo "DONE\n";
