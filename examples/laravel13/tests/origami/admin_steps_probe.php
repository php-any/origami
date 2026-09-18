<?php
/**
 * 拆开 Dashboard / Widgets HTTP 路径。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$steps = [];
$log = function ($m) use (&$steps) { echo $m."\n"; $steps[] = $m; };

$log('schema');
try {
    $page = app(\App\Filament\Pages\Dashboard::class);
    $c = $page->getWidgetsSchemaComponents($page->getWidgets());
    $log('schema_type='.gettype($c).' count='.(is_countable($c)?count($c):-1));
} catch (Throwable $e) {
    $log('schema EX: '.$e->getMessage());
}

$log('mount_stats');
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
    $log('stats_len='.strlen((string)$html));
} catch (Throwable $e) {
    $log('stats EX: '.$e->getMessage());
}

$log('mount_orders');
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\LatestOrdersWidget::class);
    $log('orders_len='.strlen((string)$html));
} catch (Throwable $e) {
    $log('orders EX: '.$e->getMessage());
    $log($e->getFile().':'.$e->getLine());
}

$log('mount_dashboard');
try {
    $html = app('livewire')->mount(\App\Filament\Pages\Dashboard::class);
    $log('dash_len='.strlen((string)$html));
    $log('dash_has='.(str_contains((string)$html,'仪表盘')||str_contains((string)$html,'Stats')||str_contains((string)$html,'订单')?'yes':'no'));
} catch (Throwable $e) {
    $log('dash EX: '.$e->getMessage());
    $log($e->getFile().':'.$e->getLine());
}

$log('http_admin');
try {
    $kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
    $req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin', 'GET');
    $req->setLaravelSession($app['session']->driver());
    $app['session']->start();
    $resp = $kernel->handle($req);
    $log('status='.$resp->getStatusCode());
    $body = (string)$resp->getContent();
    $log('len='.strlen($body));
    $log('markers='.(str_contains($body,'仪表盘')||str_contains($body,'Origami Admin')?'yes':'no'));
    $log('head='.substr(preg_replace('/\s+/',' ',strip_tags($body)),0,200));
} catch (Throwable $e) {
    $log('http EX: '.$e->getMessage());
    $log($e->getFile().':'.$e->getLine());
}

$log('DONE');
