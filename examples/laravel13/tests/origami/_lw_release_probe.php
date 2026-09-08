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
$rawAttr = $s[1] ?? '';
$snapshot = html_entity_decode($rawAttr, ENT_QUOTES);
echo "raw_attr_len=".strlen($rawAttr)." decoded_len=".strlen($snapshot)."\n";
$decoded = json_decode($snapshot, true);
echo "snap_json_ok=". (is_array($decoded)?'yes':'no')." err=".json_last_error_msg()."\n";
if (is_array($decoded)) {
    echo "keys=".json_encode(array_keys($decoded))."\n";
    echo "memo_keys=".json_encode(array_keys($decoded['memo'] ?? []))."\n";
    echo "release=".($decoded['memo']['release'] ?? $decoded['checksum'] ?? 'n/a')."\n";
    echo "snap_preview=".substr($snapshot, 0, 200)."\n";
}

$current = \Livewire\Features\SupportReleaseTokens\ReleaseToken::generate();
echo "current_release=$current\n";
echo "config_token=".config('livewire.token', 'n/a')."\n";

if (is_array($decoded)) {
    try {
        \Livewire\Features\SupportReleaseTokens\ReleaseToken::verify($decoded);
        echo "verify=ok\n";
    } catch (Throwable $e) {
        echo "verify_err=".$e->getMessage()."\n";
    }
}

$kernel->terminate($getReq, $res);

// After terminate / second request context
$current2 = \Livewire\Features\SupportReleaseTokens\ReleaseToken::generate();
echo "after_term_release=$current2 same=".($current===$current2?'yes':'no')."\n";
if (is_array($decoded)) {
    try {
        \Livewire\Features\SupportReleaseTokens\ReleaseToken::verify($decoded);
        echo "verify_after=ok\n";
    } catch (Throwable $e) {
        echo "verify_after_err=".$e->getMessage()."\n";
    }
}
