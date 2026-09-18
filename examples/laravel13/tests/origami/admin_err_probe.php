<?php
/**
 * Check if dropdown trigger slot is bound on admin HTML.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
$req = Illuminate\Http::Request::create('http://127.0.0.1:8000/admin', 'GET');
$req->setLaravelSession($app['session']->driver());
$app['session']->start();

ob_start();
try {
    $resp = $kernel->handle($req);
    $body = (string) $resp->getContent();
    echo "status=".$resp->getStatusCode()." len=".strlen($body)."\n";
    echo "has_dropdown_trigger=".(str_contains($body, 'fi-dropdown-trigger') ? 'yes' : 'no')."\n";
    echo "has_user_menu=".(str_contains($body, 'fi-user-menu') || str_contains($body, 'fi-avatar') ? 'yes' : 'no')."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}
$buf = ob_get_clean();
echo $buf;
echo "DONE\n";
