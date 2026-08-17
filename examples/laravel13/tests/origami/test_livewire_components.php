<?php
require __DIR__.'/../../vendor/autoload.php';

$checks = 0;

// Verify Livewire components extend Livewire\Component
$components = [
    'App\\Livewire\\Admin\\Login' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Dashboard' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Admins\\Index' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Admins\\Form' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Users\\Index' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Users\\Form' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Roles\\Index' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Roles\\Form' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Permissions\\Index' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Permissions\\Form' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Products\\Index' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Products\\Form' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Orders\\Index' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Orders\\Detail' => 'Livewire\\Component',
    'App\\Livewire\\Admin\\Profile' => 'Livewire\\Component',
];

foreach ($components as $comp => $parent) {
    if (!class_exists($comp)) {
        echo "FAIL: $comp does not exist\n";
        exit(1);
    }
    $reflection = new ReflectionClass($comp);
    if (!$reflection->isSubclassOf($parent)) {
        echo "FAIL: $comp is not a subclass of $parent\n";
        exit(1);
    }
    if (!$reflection->hasMethod('render')) {
        echo "FAIL: $comp has no render() method\n";
        exit(1);
    }
    $checks++;
}

echo "OK: Livewire components verified ($checks checks)\n";
