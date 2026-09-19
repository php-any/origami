<?php
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/';
$_SERVER['REQUEST_METHOD'] = 'GET';
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$rows = Illuminate\Support\Facades\DB::table('sessions')->orderByDesc('last_activity')->limit(5)->get();
echo 'session_rows='.count($rows)."\n";
foreach ($rows as $i => $r) {
    $raw = base64_decode((string) $r->payload);
    echo "row$i id=".substr((string) $r->id, 0, 16).' payload_len='.strlen((string) $r->payload).' has_login='.(str_contains($raw, 'login_admin') ? 'yes' : 'no')."\n";
    echo 'row'.$i.'_snip='.substr(preg_replace('/\s+/', ' ', $raw), 0, 220)."\n";
}
