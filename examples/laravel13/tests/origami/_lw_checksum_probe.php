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

function p($m) { file_put_contents(__DIR__.'/_login_progress.txt', $m."\n", FILE_APPEND); }
@unlink(__DIR__.'/_login_progress.txt');

$kernel = $app->make(HttpKernelContract::class);
$res = $kernel->handle(Request::create('/login', 'GET'));
$html = (string) $res->getContent();
preg_match('/wire:snapshot="([^"]+)"/', $html, $s);
$raw = html_entity_decode($s[1] ?? '', ENT_QUOTES);
p('raw_len='.strlen($raw));
file_put_contents(__DIR__.'/_snap_raw.json', $raw);

$snap = json_decode($raw, associative: true);
$check = $snap['checksum'];
unset($snap['checksum']);
$re = json_encode($snap);
p('re_len='.strlen($re));
file_put_contents(__DIR__.'/_snap_re.json', $re);

// Compare canonical payload used for checksum: original without checksum key
$rawObj = json_decode($raw, associative: true);
unset($rawObj['checksum']);
$fromRaw = json_encode($rawObj);
p('same_re='.($fromRaw === $re ? 'yes' : 'no'));
p('gen='.\Livewire\Mechanisms\HandleComponents\Checksum::generate($snap));
p('old='.$check);

// Show first difference between hashing inputs
$a = $fromRaw;
// What did generate use? Same as json_encode($snap) after unset
p('payload_head='.substr($re, 0, 200));
