<?php

use Illuminate\Contracts\Http\Kernel as HttpKernelContract;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

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
$sid = session()->getId();
$token = session()->token();
echo "saved_token=$token\n";
$kernel->terminate($getReq, $res);

$row = DB::table('sessions')->where('id', $sid)->first();
echo "row_type=".gettype($row)."\n";
if (is_object($row)) {
    echo "row_class=".get_class($row)."\n";
    echo "isset_payload=".var_export(isset($row->payload), true)."\n";
    echo "payload_preview=".substr((string)$row->payload, 0, 80)."\n";
    echo "last_activity=".var_export($row->last_activity ?? null, true)." type=".gettype($row->last_activity ?? null)."\n";
} elseif (is_array($row)) {
    echo "keys=".implode(',', array_keys($row))."\n";
    echo "payload=".substr((string)($row['payload'] ?? ''), 0, 80)."\n";
}

$decoded = base64_decode($row->payload ?? ($row['payload'] ?? ''));
echo "decoded_len=".strlen($decoded)."\n";
echo "decoded_has_token=". (str_contains($decoded, '_token') ? 'yes' : 'no')."\n";
$un = @unserialize($decoded);
echo "unserialize_ok=". (is_array($un) ? 'yes' : 'no')."\n";
if (is_array($un)) {
    echo "un_token=".($un['_token'] ?? 'MISSING')."\n";
}

$handler = app('session')->driver()->getHandler();
echo "handler=".get_class($handler)."\n";
$read = $handler->read($sid);
echo "handler_read_len=".strlen((string)$read)."\n";
echo "handler_read_has_token=". (str_contains((string)$read, '_token') ? 'yes' : 'no')."\n";
echo "handler_read_eq_decoded=". (($read === $decoded) ? 'yes' : 'no')."\n";

// expired check manually
$now = \Illuminate\Support\Carbon::now()->subMinutes(config('session.lifetime'))->getTimestamp();
$la = is_object($row) ? $row->last_activity : ($row['last_activity'] ?? null);
echo "expired_threshold=$now last_activity=$la expired=". (($la !== null && $la < $now) ? 'yes' : 'no')."\n";

// find via query
$found = DB::table('sessions')->where('id', $sid)->first();
$asObj = (object) $found;
echo "cast_isset_payload=".var_export(isset($asObj->payload), true)."\n";
echo "cast_payload_type=".gettype($asObj->payload ?? null)."\n";
if (is_object($found)) {
    // property_exists
    echo "prop_exists_payload=".var_export(property_exists($found, 'payload'), true)."\n";
}
