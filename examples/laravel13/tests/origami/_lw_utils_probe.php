<?php

require __DIR__.'/../../vendor/autoload.php';

use Livewire\Drawer\Utils;

try {
    echo Utils::stringifyHtmlAttributes(['wire:id' => 'abc', 'wire:name' => 'admin.login'])."\n";
} catch (Throwable $e) {
    echo "ERR ".$e->getMessage()."\n";
    echo "file=".$e->getFile().":".$e->getLine()."\n";
}
