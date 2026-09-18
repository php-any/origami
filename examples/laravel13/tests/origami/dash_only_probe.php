<?php
/**
 * Dashboard mount 时抓 NullValue spread / Collection::all。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

try {
    $html = app('livewire')->mount(\App\Filament\Pages\Dashboard::class);
    echo "dash_ok len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    echo "dash EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
    $p = $e->getPrevious();
    $i = 0;
    while ($p && $i < 5) {
        echo "prev$i: ".$p->getMessage()." @ ".$p->getFile().":".$p->getLine()."\n";
        $p = $p->getPrevious();
        $i++;
    }
}
echo "DONE\n";
