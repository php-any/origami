<?php
/**
 * 模拟 Livewire call 钩子：Responsable -> toResponse -> redirect effect。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Facades\Filament;
use Filament\Auth\Pages\Login;
use Livewire\Mechanisms\ExtendBlade\ExtendBlade;
use Livewire\Features\SupportRedirects\SupportRedirects;

Filament::setCurrentPanel(Filament::getPanel('admin'));

$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'POST', [], [], [], [
    'CONTENT_TYPE' => 'application/json',
    'HTTP_X_LIVEWIRE' => 'true',
]);
$app->instance('request', $req);

try {
    $html = app('livewire')->mount(Login::class);
    echo "mount_len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    echo "mount EX: ".$e->getMessage()."\n";
}

// 重新拿组件实例更难；直接走 authenticate + toResponse
$page = app(Login::class);
if (method_exists($page, 'mount')) {
    try { $page->mount(); } catch (Throwable $e) { echo "mount2 EX: ".$e->getMessage()."\n"; }
}
$page->form->fill(['email'=>'admin@example.com','password'=>'password','remember'=>false]);

// 绑定 Livewire redirector（需要 component）
$hook = new SupportRedirects();
// 手动：用 component 上的 redirect
try {
    $resp = $page->authenticate();
    echo "auth=".get_class($resp)."\n";
    echo "is_responsable=".($resp instanceof Illuminate\Contracts\Support\Responsable ? 'yes':'no')."\n";

    // 无 Livewire redirect 绑定时
    $out = $resp->toResponse(request());
    echo "toResponse_class=".get_class($out)."\n";
    if (method_exists($out, 'getTargetUrl')) {
        echo "target=".$out->getTargetUrl()."\n";
    }
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
