<?php
/**
 * 嵌套 Livewire：先 mount Login（或渲染含 @livewire Notifications 的布局），看栈。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

@unlink(storage_path('framework/stack-probe.txt'));

Livewire\before('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "BEFORE target={$cls} current={$curCls}\n",
        FILE_APPEND
    );
});

Livewire\on('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
    $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
    // 在视图编译后、真正渲染前：ExtendBlade on 应已 push
    file_put_contents(
        storage_path('framework/stack-probe.txt'),
        "ON target={$cls} current={$curCls}\n",
        FILE_APPEND
    );
    return function ($html) use ($cls) {
        $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
        $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
        file_put_contents(
            storage_path('framework/stack-probe.txt'),
            "FINISH-of-ON targetWas={$cls} current={$curCls}\n",
            FILE_APPEND
        );
        return $html;
    };
});

echo "mount_login...\n";
try {
    // Filament Login 完整页会带 layout + Notifications
    $html = app('livewire')->mount(\Filament\Auth\Pages\Login::class);
    echo "login_len=".strlen((string)$html)."\n";
    echo "has_err=".(str_contains((string)$html, 'getBroadcastChannel') ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "---stack---\n";
echo @file_get_contents(storage_path('framework/stack-probe.txt'));
echo "DONE\n";
