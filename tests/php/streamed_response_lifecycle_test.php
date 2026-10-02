<?php
use Symfony\Component\HttpFoundation\StreamedResponse;

$calls = 0;
$response = new StreamedResponse(function () use (&$calls) {
    $calls++;
    echo 'stream';
});
ob_start();
$response->sendContent();
$response->sendContent();
$body = ob_get_clean();
if ($body !== 'stream' || $calls !== 1 || $response->getContent() !== false) {
    throw new RuntimeException('StreamedResponse callback/idempotence failed');
}
$chunks = new StreamedResponse(['a', 'b']);
ob_start();
ob_start();
$chunks->sendContent();
ob_end_flush();
$body = ob_get_clean();
if ($body !== 'ab') { throw new RuntimeException('chunks failed'); }
$source = (function () { yield 'x'; yield 'y'; })();
if (!$source->valid() || $source->current() !== 'x') { throw new RuntimeException('generator did not initialize'); }
$generator = new StreamedResponse($source);
ob_start();
ob_start();
$generator->sendContent();
ob_end_flush();
$generated = ob_get_clean();
if ($generated !== 'xy') { throw new RuntimeException('generator chunks failed: '.$generated); }
echo "OK: streamed response lifecycle\n";
