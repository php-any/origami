<?php
/**
 * 渲染官方 filament icon-button 视图，检查 class 是否变成 0=。
 */
chdir(__DIR__.'/../..');
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Support\Icons\Heroicon;
use Illuminate\View\ComponentAttributeBag;

$html = view('filament::components.icon-button', [
    'label' => 'Expand sidebar',
    'icon' => Heroicon::OutlinedBars3,
    'iconAlias' => null,
    'iconSize' => 'lg',
    'color' => 'gray',
    'tag' => 'button',
    'type' => 'button',
    'disabled' => false,
    'tooltip' => null,
    'attributes' => new ComponentAttributeBag([
        'x-cloak' => true,
        'x-data' => '{}',
        'aria-controls' => 'fi-main-sidebar',
        'x-bind:aria-expanded' => '$store.sidebar.isOpen',
        'x-on:click' => '$store.sidebar.open()',
        'x-show' => '! $store.sidebar.isOpen',
        'class' => 'fi-topbar-open-sidebar-btn',
    ]),
])->render();

$dbg = __DIR__.'/../../storage/origami-debug';
if (!is_dir($dbg)) {
    mkdir($dbg, 0755, true);
}
file_put_contents($dbg.'/icon-button-view.html', $html);
echo 'len='.strlen($html)."\n";
echo 'has_numeric='.(preg_match('/\s0="/', $html) ? 'yes' : 'no')."\n";
echo 'has_class_fi='.(str_contains($html, 'fi-icon-btn') ? 'yes' : 'no')."\n";
echo 'has_aria='.(str_contains($html, 'aria-label') ? 'yes' : 'no')."\n";
echo substr(preg_replace('/\s+/', ' ', $html), 0, 500), "\n";
