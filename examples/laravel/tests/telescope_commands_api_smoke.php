<?php

/**
 * CommandsController@index 端到端验收（Illuminate Request + 容器仓储）。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Illuminate\Http\Request;
use Laravel\Telescope\Http\Controllers\CommandsController;
use Symfony\Component\HttpFoundation\Request as SymfonyRequest;

$request = Request::createFromBase(SymfonyRequest::create(
    '/telescope/telescope-api/commands',
    'POST',
    ['tag' => '', 'before' => '', 'take' => 50, 'family_hash' => '']
));

$storage = telescope_entries_repository();
$controller = new CommandsController();
$response = $controller->index($request, $storage);

$content = $response->getContent();
$data = json_decode($content, true);
if (!is_array($data) || !array_key_exists('entries', $data)) {
    echo "FAIL: unexpected response: $content\n";
    exit(1);
}

echo "PASS\n";
