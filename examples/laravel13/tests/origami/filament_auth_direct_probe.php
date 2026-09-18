<?php
/**
 * 直接调用 Filament Login::authenticate（绕过 Livewire update），定位卡点。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Filament\Facades\Filament;
use Filament\Auth\Pages\Login;
use Illuminate\Support\Facades\Hash;

echo "admin_count=".Admin::count()."\n";
$admin = Admin::where('email', 'admin@example.com')->first();
if (!$admin) {
    echo "FAIL: no admin seed\n";
    exit(1);
}
echo "admin_id=".$admin->id." hash_ok=".(Hash::check('password', $admin->password) ? 'yes' : 'no')."\n";
echo "canAccess=".(method_exists($admin, 'canAccessPanel') ? ($admin->canAccessPanel(Filament::getCurrentOrDefaultPanel() ?? Filament::getPanel('admin')) ? 'yes' : 'no') : 'n/a')."\n";

// 设置当前请求，避免 session 问题
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin/login', 'GET');
$app->instance('request', $req);
Illuminate\Support\Facades\Facade::clearResolvedInstance('request');

try {
    Filament::setCurrentPanel(Filament::getPanel('admin'));
} catch (Throwable $e) {
    echo "setPanel EX: ".$e->getMessage()."\n";
}

echo "panel=".optional(Filament::getCurrentOrDefaultPanel())->getId()."\n";

try {
    /** @var Login $page */
    $page = app(Login::class);
    // Livewire mount 简化：填表
    if (method_exists($page, 'mount')) {
        $page->mount();
    }
    $page->form->fill([
        'email' => 'admin@example.com',
        'password' => 'password',
        'remember' => false,
    ]);
    echo "form_state=".json_encode($page->form->getState())."\n";
    $resp = $page->authenticate();
    echo "auth_result=".($resp === null ? 'null' : get_class($resp))."\n";
    echo "guard_check=".(Filament::auth()->check() ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
    echo "trace:\n";
    foreach (array_slice($e->getTrace(), 0, 8) as $i => $t) {
        echo "#$i ".($t['class'] ?? '').($t['type'] ?? '').($t['function'] ?? '')." @ ".($t['file'] ?? '').":".($t['line'] ?? '')."\n";
    }
}
echo "DONE\n";
