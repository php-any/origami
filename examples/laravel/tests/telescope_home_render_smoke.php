<?php

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

$controller = new \Laravel\Telescope\Http\Controllers\HomeController();
try {
    $result = $controller->index();
    echo 'result class: ' . (is_object($result) ? get_class($result) : gettype($result)) . "\n";
    if (is_object($result) && method_exists($result, 'render')) {
        $html = $result->render();
        echo 'len=' . strlen($html) . "\n";
        echo str_contains($html, 'id="telescope"') ? "PASS\n" : "FAIL\n";
    }
} catch (Throwable $e) {
    echo 'ERR: ' . $e->getMessage() . "\n";
}
