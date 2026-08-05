<?php

/**
 * RequestsController@show 官方路由存在。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

$found = false;
foreach (illuminate_router()->getRoutes() as $route) {
    if ((string) $route->uri() === 'telescope/telescope-api/requests/{telescopeEntryId}'
        && str_ends_with((string) $route->getActionName(), 'RequestsController@show')) {
        $found = true;
        break;
    }
}
if (!$found) {
    echo "FAIL: RequestsController@show route missing\n";
    exit(1);
}

echo "PASS\n";
