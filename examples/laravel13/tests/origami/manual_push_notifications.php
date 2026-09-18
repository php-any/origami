<?php
/**
 * 在 Notifications render 的 before 阶段手动调用 ExtendBlade 实例的 startLivewireRendering。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

@unlink(storage_path('framework/manual-push2.txt'));

// 找到已 boot 的 ExtendBlade 机制实例
$eb = null;
$manager = app('livewire');
// LivewireManager 不一定暴露 mechanisms；直接 new 并只借方法会用不同 static！
// 必须调用「同一个类」上的实例方法，使 static::$livewireComponents 为同一槽。

Livewire\before('render', function ($target, $view) {
    $cls = is_object($target) ? get_class($target) : gettype($target);
    if (str_contains($cls, 'Notifications')) {
        // 通过匿名子类调用父类 protected static？改为：复制 start 逻辑不可行。
        // 使用 invade 调已有 listener？直接 new ExtendBlade 并 start —— PHP 中 static 属性按类声明处共享。
        $inst = new Livewire\Mechanisms\ExtendBlade\ExtendBlade();
        $inst->startLivewireRendering($target);
        $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
        $curCls = is_object($cur) ? get_class($cur) : var_export($cur, true);
        file_put_contents(
            storage_path('framework/manual-push2.txt'),
            "manual_push target={$cls} currentAfter={$curCls}\n",
            FILE_APPEND
        );
    }
});

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$resp = $http->handle($req);
$c = (string)$resp->getContent();
echo "status=".$resp->getStatusCode()."\n";
echo "has_err=".(str_contains($c, 'getBroadcastChannel') || str_contains($c, 'Login::getBroadcastChannel') ? 'yes' : 'no')."\n";
// 也检查 stdout 异常是否已进 content — probe 里异常在 stdout
echo "---manual---\n";
echo @file_get_contents(storage_path('framework/manual-push2.txt'));
echo "DONE\n";
