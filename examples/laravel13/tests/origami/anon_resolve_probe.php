<?php
/**
 * AnonymousComponent resolve must keep view name; factory->exists must find filament view.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$viewName = 'filament-panels::components.page.simple';
echo "exists=".(app('view')->exists($viewName)?'yes':'no')."\n";

$component = Illuminate\View\AnonymousComponent::resolve([
    'view' => $viewName,
    'data' => [],
]);
echo "class=".get_class($component)."\n";

// Invade protected view
$ref = new ReflectionProperty($component, 'view');
$ref->setAccessible(true);
echo "view_prop=".json_encode($ref->getValue($component))."\n";

$ref2 = new ReflectionProperty($component, 'data');
$ref2->setAccessible(true);
echo "data_prop=".json_encode($ref2->getValue($component))."\n";

echo "render=".json_encode($component->render())."\n";
echo "resolveView=".json_encode($component->resolveView())."\n";
