<?php
class DelegateInner extends ArrayIterator {
    public function extra($prefix) { return $prefix . $this->current(); }
    public function fail() { throw new RuntimeException('inner failure'); }
}
class DelegateFilter extends FilterIterator { public function accept(): bool { return true; } }
$inner = new DelegateInner(['value']);
$outer = new DelegateFilter(new IteratorIterator($inner));
$outer->rewind();
if ($outer->extra('prefix:') !== 'prefix:value' || !method_exists($outer, 'extra') || method_exists($outer, '__call')) { throw new Exception('Iterator delegation failed'); }
try { $outer->fail(); throw new Exception('Exception swallowed'); } catch (RuntimeException $e) { if ($e->getMessage() !== 'inner failure') { throw $e; } }
echo "iterator method delegation passed\n";
