<?php
/**
 * Schema::toEmbeddedHtml 返回类型。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Filament\Schemas\Schema;
use Filament\Schemas\Components\Component;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

// 直接复现 toEmbeddedHtml 写法
$dummy = new class {
    public function renderEmbeddedHtml(): string { return '<div>ok</div>'; }
    public function toEmbeddedHtml(): string {
        return Component::withVisibilityCache(fn (): string => $this->renderEmbeddedHtml());
    }
};
try {
    $r = $dummy->toEmbeddedHtml();
    echo "dummy_type=".gettype($r)." val=".var_export($r,true)."\n";
} catch (Throwable $e) {
    echo "dummy EX: ".$e->getMessage()."\n";
}

// Stats widget mount 内部会走到 schema toEmbeddedHtml
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
    echo "stats_len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    echo "stats EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
