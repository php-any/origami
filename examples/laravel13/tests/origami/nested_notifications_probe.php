<?php
/**
 * Nest Notifications inside a Livewire render to see $this class for getBroadcastChannel.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

try {
    $html = Illuminate\Support\Facades\Blade::render(
        '@livewire(\\Filament\\Notifications\\Livewire\\Notifications::class)'
    );
    echo "blade_livewire_ok len=".strlen($html)."\n";
    echo (str_contains($html, 'fi-no') || str_contains($html, 'wire:id')) ? "markers:yes\n" : "markers:no\n";
} catch (Throwable $e) {
    echo "blade_EX: ".$e->getMessage()."\n";
}

// Nested: render Login page fragment that includes notifications via layout
try {
    $html2 = app('livewire')->mount(\Filament\Auth\Pages\Login::class);
    echo "login_mount_ok type=".gettype($html2)." len=".(is_string($html2)?strlen($html2):0)."\n";
    if (is_string($html2) && str_contains($html2, 'getBroadcastChannel')) {
        echo "login_html_has_error_text=yes\n";
    }
} catch (Throwable $e) {
    echo "login_EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
