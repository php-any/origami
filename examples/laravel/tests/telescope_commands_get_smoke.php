<?php
require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Illuminate\Http\Request;
use Laravel\Telescope\EntryType;
use Laravel\Telescope\Storage\EntryQueryOptions;
use Symfony\Component\HttpFoundation\Request as SymfonyRequest;

$request = Request::createFromBase(SymfonyRequest::create('/x', 'POST', [
    'tag' => '', 'before' => '', 'take' => 50, 'family_hash' => '',
]));
$repo = telescope_entries_repository();
$options = EntryQueryOptions::fromRequest($request);
$entries = $repo->get(EntryType::COMMAND, $options);
echo 'count=' . count($entries) . "\nPASS\n";
