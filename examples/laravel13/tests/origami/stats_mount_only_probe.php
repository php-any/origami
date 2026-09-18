<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "boot_ok\n";

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));
echo "auth_ok\n";

$html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
echo "stats_len=".strlen((string)$html)."\n";
echo "DONE\n";
