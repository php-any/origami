<?php
require __DIR__ . '/../../vendor/autoload.php';

use Livewire\Features\SupportPageComponents\PageComponentConfig;
use Livewire\Mechanisms\HandleComponents\ViewContext;

$vc = new ViewContext();
echo "ViewContext ok\n";

$c = new PageComponentConfig('component', 'layouts.app', 'slot', []);
echo "PageComponentConfig ok viewContext=";
echo $c->viewContext === null ? 'NULL' : get_class($c->viewContext);
echo "\n";
