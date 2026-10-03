<?php
class NativeThrowParent extends Exception {}
class NativeThrowChild extends NativeThrowParent {}
function throw_nominal_check(Throwable $error): bool { return $error instanceof NativeThrowParent; }
try { throw new NativeThrowChild('caught'); } catch (NativeThrowParent $error) {
    if (!throw_nominal_check($error) || $error->getMessage() !== 'caught') throw new Exception('throw declaration lost');
}
echo "throw nominal contract OK\n";
