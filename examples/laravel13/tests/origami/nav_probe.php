<?php
/**
 * Narrow getNavigation recursion.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$panel = Filament::getPanel('admin');
echo "panel_ok\n";

try {
    $nav = $panel->getNavigation();
    echo "nav_type=".gettype($nav);
    if (is_object($nav)) {
        echo " class=".get_class($nav);
    }
    echo "\n";
    if (is_object($nav) && method_exists($nav, 'all')) {
        $a = $nav->all();
        echo "nav_all=".gettype($a)."\n";
    }
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo "FILE: ".$e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
