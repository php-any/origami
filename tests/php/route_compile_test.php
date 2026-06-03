<?php
require __DIR__ . '/../../examples/laravel/vendor/autoload.php';

$route = new \Symfony\Component\Routing\Route('/');
try {
    $compiled = $route->compile();
    echo "regex=" . $compiled->getRegex() . "\n";
} catch (\Throwable $e) {
    echo "compile error: " . $e->getMessage() . "\n";
}

$route2 = new \Symfony\Component\Routing\Route('/users/{id}');
try {
    $compiled2 = $route2->compile();
    echo "regex2=" . $compiled2->getRegex() . "\n";
} catch (\Throwable $e) {
    echo "compile2 error: " . $e->getMessage() . "\n";
}
