<?php
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';

$r = Illuminate\Http\Request::capture();
echo "pathInfo=" . $r->getPathInfo() . "\n";
echo "requestUri=" . $r->getRequestUri() . "\n";
echo "SCRIPT_NAME=" . ($r->server->get('SCRIPT_NAME') ?? '') . "\n";
