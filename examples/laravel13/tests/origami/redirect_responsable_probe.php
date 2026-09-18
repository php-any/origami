<?php
/**
 * 验证 LoginResponse instanceof Responsable，以及 Livewire Redirector 绑定。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Facades\Filament;
use Filament\Auth\Pages\Login;
use Filament\Auth\Http\Responses\LoginResponse;
use Illuminate\Contracts\Support\Responsable;
use Livewire\Features\SupportRedirects\SupportRedirects;
use Livewire\Features\SupportRedirects\Redirector as LwRedirector;

Filament::setCurrentPanel(Filament::getPanel('admin'));

$resp = app(LoginResponse::class);
echo "class=".get_class($resp)."\n";
echo "instanceof_Responsable=".($resp instanceof Responsable ? 'yes':'no')."\n";
echo "instanceof_FilamentContract=".($resp instanceof \Filament\Auth\Http\Responses\Contracts\LoginResponse ? 'yes':'no')."\n";
echo "implements=".implode(',', class_implements($resp) ?: [])."\n";

// 模拟 Livewire：挂 Redirector 再 toResponse
$page = app(Login::class);
app()->bind('redirect', function () use ($page) {
    $r = app(LwRedirector::class)->component($page);
    if (app()->has('session.store')) {
        $r->setSession(app('session.store'));
    }
    return $r;
});

echo "redirect_bound=".get_class(app('redirect'))."\n";

try {
    $page->form->fill(['email'=>'admin@example.com','password'=>'password','remember'=>false]);
    $auth = $page->authenticate();
    echo "auth_class=".get_class($auth)."\n";
    echo "auth_responsable=".($auth instanceof Responsable ? 'yes':'no')."\n";

    if ($auth instanceof Responsable) {
        $out = $auth->toResponse(request());
        echo "toResponse_class=".get_class($out)."\n";
    }

    $store = \Livewire\store($page);
    $to = $store->get('redirect');
    echo "store_redirect=".var_export($to, true)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
