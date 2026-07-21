<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Container\Container;
use Illuminate\Translation\ArrayLoader;
use Illuminate\Translation\Translator;

$container = new Container();
$container->singleton('translator', function () {
    return new Translator(new ArrayLoader(), 'en');
});

$t = $container->make('translator');
echo $t->get('validation.required', ['attribute' => 'email']) . "\n";
