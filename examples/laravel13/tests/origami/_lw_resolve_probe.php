<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;

require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

$kernel = $app->make(HttpKernelContract::class);
$getReq = Request::create('/login', 'GET');
$res = $kernel->handle($getReq);
$html = (string) $res->getContent();
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
$snapshot = html_entity_decode($s[1] ?? '', ENT_QUOTES);
$decoded = json_decode($snapshot, true);
echo "name=".($decoded['memo']['name'] ?? '')."\n";
echo "release=".($decoded['memo']['release'] ?? '')."\n";
echo "checksum=".($decoded['checksum'] ?? '')."\n";
echo "full_snap_len=".strlen($snapshot)."\n";

try {
    $class = app('livewire.factory')->resolveComponentClass($decoded['memo']['name']);
    echo "resolved=$class\n";
    $gen = \Livewire\Features\SupportReleaseTokens\ReleaseToken::generate($class);
    echo "generated=$gen match=".(($decoded['memo']['release'] ?? '') === $gen ? 'yes' : 'no')."\n";
    \Livewire\Features\SupportReleaseTokens\ReleaseToken::verify($decoded);
    echo "verify=ok\n";
} catch (Throwable $e) {
    echo "ERR ".get_class($e)." ".$e->getMessage()."\n";
}

// Compare with how Livewire would decode from request components payload
$payloadSnap = json_decode(json_encode(['snapshot' => $snapshot]), true);
$inner = json_decode($payloadSnap['snapshot'], associative: true);
echo "roundtrip_name=".($inner['memo']['name'] ?? 'MISSING')."\n";
echo "roundtrip_ok=". (is_array($inner)?'yes':'no')."\n";
