<?php

require __DIR__ . '/../../vendor/autoload.php';

$c = new Illuminate\Container\Container();
$c->bind('foo', function ($app) {
    return 'bar';
});
$v = $c->make('foo');
if ($v !== 'bar') {
    fwrite(STDERR, "make foo failed: " . var_export($v, true) . "\n");
    exit(1);
}
echo "container_bind_ok\n";
