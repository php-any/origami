<?php
/**
 * Dump get_defined_vars keys at the moment @capture would call it, inside logo view.
 * We compile a tiny blade that mirrors logo's capture preamble.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

// Render via component tag compiler path
try {
    $view = Illuminate\Support\Facades\View::make('filament-panels::components.logo', [
        'attributes' => new Illuminate\View\ComponentAttributeBag([]),
    ]);
    echo "view_name=".$view->name()."\n";
    $data = $view->getData();
    echo "view_data_keys=".implode(',', array_keys($data))."\n";
    echo "view_has_attr=".(array_key_exists('attributes', $data)?'yes':'no')."\n";
    if (isset($data['attributes'])) {
        echo "view_attr=".get_class($data['attributes'])."\n";
    }
    $html = $view->render();
    echo "render_ok len=".strlen($html)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

// Component resolve path
try {
    $html2 = Illuminate\Support\Facades\Blade::render(
        '@component("filament-panels::components.logo") @endcomponent'
    );
    echo "component_ok len=".strlen($html2)."\n";
} catch (Throwable $e) {
    echo "COMP_EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
