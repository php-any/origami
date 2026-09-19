<?php

/**
 * Telescope entries 多包一层：json_encode(collect(JsonSerializable[])) 与 repo->get。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

class TelJsonWrap_Box implements JsonSerializable
{
    public function __construct(private string $id)
    {
    }

    public function jsonSerialize(): mixed
    {
        return collect([
            'id' => $this->id,
            'content' => ['method' => 'GET', 'uri' => '/'],
        ])->all();
    }
}

$c = collect([new TelJsonWrap_Box('a'), new TelJsonWrap_Box('b')]);
$all = $c->all();
echo 'collect_count='.(is_countable($all) ? count($all) : -1)."\n";
echo 'collect_all='.json_encode($all)."\n";
echo 'collect_json='.json_encode(['entries' => $c])."\n";

$resp = response()->json(['entries' => $c, 'status' => 'enabled']);
echo 'jsonresp='.$resp->getContent()."\n";

$repo = $app->make(Laravel\Telescope\Contracts\EntriesRepository::class);
$opts = Laravel\Telescope\Storage\EntryQueryOptions::fromRequest(
    Illuminate\Http\Request::create('/telescope/telescope-api/requests', 'POST', ['take' => 2])
);
$got = $repo->get(Laravel\Telescope\EntryType::REQUEST, $opts);
echo 'repo_class='.get_class($got)."\n";
$repoAll = $got->all();
echo 'repo_all_count='.(is_countable($repoAll) ? count($repoAll) : -1)."\n";
echo 'repo_all0_type='.gettype($repoAll[0] ?? null)."\n";
echo 'repo_all0_class='.(is_object($repoAll[0] ?? null) ? get_class($repoAll[0]) : '-')."\n";
$repoJson = json_encode(['entries' => $got, 'status' => 'enabled']);
$decoded = json_decode($repoJson, true);
$entries = $decoded['entries'] ?? null;
echo 'entries_is_list='.(is_array($entries) && array_is_list($entries) ? '1' : '0')."\n";
echo 'entries0_is_list='.(isset($entries[0]) && is_array($entries[0]) && array_is_list($entries[0]) ? '1' : '0')."\n";
echo 'entries0_method='.($entries[0]['content']['method'] ?? '-')."\n";
echo 'entries_count='.(is_array($entries) ? count($entries) : -1)."\n";
