<?php
class PregCallbackReceiverProbe {
    public $calls = 0;
    public function replace($match) { ++$this->calls; return strtoupper($match[0]); }
}
$object = new PregCallbackReceiverProbe();
if (preg_replace_callback('/[a-z]+/', [$object, 'replace'], 'one two') !== 'ONE TWO' || $object->calls !== 2) {
    throw new \Exception('Regex callback receiver lost');
}
if (preg_replace_callback('/[a-z]+/', function ($m) { return strlen($m[0]); }, 'one two') !== '3 3') {
    throw new \Exception('Regex closure callback failed');
}
try {
    preg_replace_callback('/x/', function ($m) { throw new \RuntimeException('callback-error'); }, 'x');
    throw new \Exception('Callback exception swallowed');
} catch (\RuntimeException $e) {
    if ($e->getMessage() !== 'callback-error') { throw $e; }
}
echo "regex callback receiver passed\n";
