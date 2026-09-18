<?php
/**
 * Check instanceof Htmlable for Schema content.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$panel = Filament\Facades\Filament::getPanel('admin');
Filament\Facades\Filament::setCurrentPanel($panel);
$login = app(Filament\Auth\Pages\Login::class);
$login->mount();
$c = $login->content;
echo "class=".get_class($c)."\n";
echo "htmlable=".( ($c instanceof Illuminate\Contracts\Support\Htmlable) ? 'yes':'no')."\n";
echo "parents=".json_encode(class_parents($c))."\n";
echo "implements=".json_encode(class_implements($c))."\n";
if ($c instanceof Illuminate\Contracts\Support\Htmlable) {
    try {
        $html = $c->toHtml();
        echo "toHtml_len=".strlen((string)$html)."\n";
        echo "toHtml_head=".substr(preg_replace('/\s+/',' ',(string)$html),0,200)."\n";
    } catch (Throwable $e) {
        echo "toHtml_EX=".get_class($e).": ".$e->getMessage()."\n";
        echo "at ".$e->getFile().":".$e->getLine()."\n";
    }
}
echo "DONE\n";
