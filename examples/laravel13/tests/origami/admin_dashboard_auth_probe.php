<?php
/**
 * Auth::login 后 GET /admin，抓完整异常。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;

$admin = Admin::where('email', 'admin@example.com')->first();
Auth::guard('admin')->login($admin);
echo "logged_in=yes\n";

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http\Request::create('http://127.0.0.1:8000/admin', 'GET');
$req->setLaravelSession($app['session']->driver());
$app['session']->start();

try {
    $resp = $kernel->handle($req);
    echo "status=".$resp->getStatusCode()."\n";
    echo "loc=".(string)$resp->headers->get('Location')."\n";
    $c = (string)$resp->getContent();
    echo "len=".strlen($c)."\n";
    echo "markers=".(str_contains($c,'Origami Admin')||str_contains($c,'Dashboard')||str_contains($c,'仪表盘')?'yes':'no')."\n";
    echo "head=".substr(preg_replace('/\s+/',' ',strip_tags($c)),0,250)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
    foreach (array_slice($e->getTrace(), 0, 12) as $i => $t) {
        echo "#$i ".($t['class'] ?? '').($t['type'] ?? '').($t['function'] ?? '')." ".$e->getFile()."\n";
        echo "   @ ".($t['file'] ?? '').':'.($t['line'] ?? '')."\n";
    }
}
echo "DONE\n";
