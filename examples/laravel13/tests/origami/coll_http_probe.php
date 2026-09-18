<?php
/**
 * GET /admin 时抓 Collection::all 递归：打印 items 类型。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Illuminate\Support\Collection;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

// 在 all 外包一层诊断：用 class_alias 不现实；直接测可疑路径
$c = collect([1,2,3]);
echo "basic_all=".gettype($c->all())."\n";

$c2 = new Collection(new Collection([1,2]));
echo "nested_all=".gettype($c2->all())." count=".count($c2->all())."\n";

// 导航相关
try {
    $panel = Filament::getPanel('admin');
    $nav = $panel->getNavigation();
    echo "nav_type=".gettype($nav).(is_object($nav)?(' '.get_class($nav)):'')."\n";
    if (is_object($nav) && method_exists($nav, 'all')) {
        $a = $nav->all();
        echo "nav_all_type=".gettype($a)."\n";
    }
} catch (Throwable $e) {
    echo "nav EX: ".$e->getMessage()."\n";
}

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin', 'GET');
$req->setLaravelSession($app['session']->driver());
$app['session']->start();
try {
    $resp = $kernel->handle($req);
    echo "status=".$resp->getStatusCode()." len=".strlen((string)$resp->getContent())."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}
echo "DONE\n";
