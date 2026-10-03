<?php
class RejectingFilter extends FilterIterator {
    public function accept(): bool { throw new RuntimeException('accept'); }
}
try { (new RejectingFilter(new ArrayIterator([1])))->rewind(); throw new Exception('missing accept exception'); }
catch (RuntimeException $error) { if ($error->getMessage() !== 'accept') throw $error; }
class FailingInner extends ArrayIterator {
    public function valid(): bool { throw new RuntimeException('valid'); }
}
class AcceptingFilter extends FilterIterator { public function accept(): bool { return true; } }
try { (new AcceptingFilter(new FailingInner([1])))->rewind(); throw new Exception('missing valid exception'); }
catch (RuntimeException $error) { if ($error->getMessage() !== 'valid') throw $error; }
echo "SPL exception propagation OK\n";
