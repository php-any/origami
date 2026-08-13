<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Container\Container;

$c = new Container();

$c->bind('greeting', function () {
    return 'hello';
});

if ($c->make('greeting') !== 'hello') {
    echo "FAIL: bind/make\n";
    exit(1);
}

$c->singleton('counter', function () {
    $o = new stdClass();
    $o->n = 0;
    return $o;
});

$a = $c->make('counter');
$b = $c->make('counter');
$a->n = 7;
if ($b->n !== 7) {
    echo "FAIL: singleton not shared\n";
    exit(1);
}

$c->instance('fixed', 'origami');
if ($c->make('fixed') !== 'origami') {
    echo "FAIL: instance\n";
    exit(1);
}

if (!$c->bound('greeting') || !$c->bound('counter')) {
    echo "FAIL: bound\n";
    exit(1);
}

echo "PASS\n";
