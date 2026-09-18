<?php
/**
 * 仿 ExtendBlade：类方法内注册 on('render')，嵌套 HTTP 时记录 $this。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

class ExtendBladeThisSpy
{
    public function boot()
    {
        Livewire\on('render', function ($target, $view) {
            $thisClass = is_object($this) ? get_class($this) : var_export($this, true);
            $targetClass = is_object($target) ? get_class($target) : gettype($target);
            $cur = Livewire\Mechanisms\ExtendBlade\ExtendBlade::currentRendering();
            $curClass = is_object($cur) ? get_class($cur) : var_export($cur, true);
            file_put_contents(
                storage_path('framework/this-spy.txt'),
                "this={$thisClass} target={$targetClass} current={$curClass}\n",
                FILE_APPEND
            );
            // 尝试像 ExtendBlade 一样 push（经由公开 API 不行）；直接调 start 需实例
            return function ($html) {
                return $html;
            };
        });
    }
}

@unlink(storage_path('framework/this-spy.txt'));
(new ExtendBladeThisSpy())->boot();

$http = app(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$resp = $http->handle($req);
$c = (string)$resp->getContent();
echo "status=".$resp->getStatusCode()."\n";
echo "has_broadcast_msg=".(str_contains($c, 'getBroadcastChannel') || file_exists(storage_path('framework/this-spy.txt')) ? 'check' : 'no')."\n";
echo "---spy---\n";
echo @file_get_contents(storage_path('framework/this-spy.txt'));
echo "DONE\n";
