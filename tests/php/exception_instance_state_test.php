<?php
function expectExceptionState($condition) {
    if (!$condition) { throw new Exception('exception instance state mismatch'); }
}
class IsolatedException extends Exception {}
$first = new IsolatedException('first', 17);
$trace = $first->getTraceAsString();
$second = new IsolatedException('second', 23, $first);
expectExceptionState($first->getMessage() === 'first');
expectExceptionState($first->getCode() === 17);
expectExceptionState($second->getMessage() === 'second');
expectExceptionState($second->getCode() === 23);
expectExceptionState($second->getPrevious() === $first);
expectExceptionState($first->getTraceAsString() === $trace);
foreach ($first->getTrace() as $frame) {
    expectExceptionState(is_array($frame));
}
$error = new ErrorException('warning', 5, E_WARNING);
expectExceptionState($error->getMessage() === 'warning');
expectExceptionState($error->getSeverity() === E_WARNING);
echo "exception instance state ok\n";
