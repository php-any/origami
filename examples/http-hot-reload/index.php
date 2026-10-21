<?php

use Net\Http\Server;

require __DIR__.'/bootstrap.php';
$port = getenv('HTTP_PORT');
$port = ($port === false || $port === '') ? 18087 : (int) $port;
$server = new Server('127.0.0.1', port: $port);
require __DIR__.'/routes.php';
echo "Listening on http://127.0.0.1:{$port}\n";
$server->run();
