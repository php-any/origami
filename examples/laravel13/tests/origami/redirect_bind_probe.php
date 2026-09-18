<?php
/**
 * 验证 bind('redirect') 是否覆盖已有 instance。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Livewire\Features\SupportRedirects\Redirector as LwRedirector;
use Filament\Auth\Pages\Login;
use Filament\Facades\Filament;

Filament::setCurrentPanel(Filament::getPanel('admin'));

echo "before=".get_class(app('redirect'))."\n";
echo "has_instance=".(isset(app()->getBindings()['redirect']) || app()->bound('redirect') ? 'bound' : 'no')."\n";

$page = app(Login::class);
$before = app('redirect');

app()->bind('redirect', function () use ($page) {
    echo "factory_called\n";
    return app(LwRedirector::class)->component($page);
});

echo "after_bind_class=".get_class(app('redirect'))."\n";
echo "same_as_before=".($before === app('redirect') ? 'yes':'no')."\n";

// instance 再 bind
app()->instance('redirect', $before);
echo "after_instance=".get_class(app('redirect'))."\n";
app()->bind('redirect', function () use ($page) {
    echo "factory2_called\n";
    return app(LwRedirector::class)->component($page);
});
echo "after_rebind=".get_class(app('redirect'))."\n";

echo "DONE\n";
