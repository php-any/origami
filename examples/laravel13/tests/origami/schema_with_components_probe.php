<?php
/**
 * 带组件的 Schema::toEmbeddedHtml。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Filament\Schemas\Schema;
use Filament\Schemas\Components\Section;
use Filament\Schemas\Components\Text;
use Filament\Schemas\Components\Component;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

fwrite(STDERR, "step1\n");

// 最小带组件 schema（需 livewire 宿主？）
try {
    $widget = app(\App\Filament\Widgets\StatsOverviewWidget::class);
    // 走 Livewire 的 mount 初始化
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
    fwrite(STDERR, "mount_ok len=".strlen((string)$html)."\n");
    echo "mount_len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    fwrite(STDERR, "mount_fail\n");
    echo "mount EX: ".$e->getMessage()."\n";
    echo "file=".$e->getFile().":".$e->getLine()."\n";
    echo "type=".get_class($e)."\n";
}

fwrite(STDERR, "step2\n");

// 手工构造简单 schema + Text 组件
try {
    $schema = Schema::make()
        ->components([
            Text::make('hello'),
        ]);
    // Schema 可能需要 livewire
    if (method_exists($schema, 'livewire')) {
        $lw = app(\App\Filament\Widgets\StatsOverviewWidget::class);
        $schema->livewire($lw);
    }
    echo "components=".count($schema->getComponents(withHidden: true))."\n";

    $ref = new ReflectionClass($schema);
    $m = $ref->getMethod('renderEmbeddedHtml');
    $m->setAccessible(true);
    try {
        $html = $m->invoke($schema);
        echo "render_type=".gettype($html)."\n";
        if (is_string($html)) echo "render_len=".strlen($html)."\n";
        if (is_object($html)) echo "render_class=".get_class($html)."\n";
    } catch (Throwable $e) {
        echo "render EX: ".$e->getMessage()."\n";
        echo $e->getFile().":".$e->getLine()."\n";
    }

    try {
        $html = $schema->toEmbeddedHtml();
        echo "embed_type=".gettype($html)."\n";
        if (is_string($html)) echo "embed_len=".strlen($html)."\n";
        if (is_object($html)) echo "embed_class=".get_class($html)."\n";
    } catch (Throwable $e) {
        echo "embed EX: ".$e->getMessage()."\n";
        echo $e->getFile().":".$e->getLine()."\n";
    }
} catch (Throwable $e) {
    echo "build EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

fwrite(STDERR, "done\n");
echo "DONE\n";
