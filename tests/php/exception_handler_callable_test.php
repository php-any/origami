<?php
function handlerCallableCheck($ok, $message) {
    if (!$ok) throw new Exception($message);
}
function namedExceptionHandler(Throwable $e) { echo 'named:' . $e->getMessage() . "\n"; }
class HandlerCallableFixture {
    public function instanceHandler(Throwable $e) { echo 'instance:' . $e->getMessage() . "\n"; }
    public static function staticHandler(Throwable ...$errors) { echo 'static:' . count($errors) . ':' . $errors[0]->getMessage() . "\n"; }
    private function privateHandler(Throwable $e) { echo 'private:' . $e->getMessage() . "\n"; }
    public function __invoke(Throwable $e) { echo 'invokable:' . $e->getMessage() . "\n"; }
}
$object = new HandlerCallableFixture();
$callbacks = ['namedExceptionHandler', [$object, 'instanceHandler'], ['HandlerCallableFixture', 'staticHandler'], 'HandlerCallableFixture::staticHandler', $object];
$previous = null;
foreach ($callbacks as $callback) {
    handlerCallableCheck(set_exception_handler($callback) === $previous, 'original callback identity');
    $previous = $callback;
}
foreach ([['HandlerCallableFixture', 'instanceHandler'], [$object, 'privateHandler'], [1 => 'staticHandler', 0 => 1]] as $bad) {
    try {
        set_exception_handler($bad);
        throw new Exception('accepted inaccessible callback');
    } catch (TypeError $e) {}
}
register_shutdown_function(function () { echo "shutdown\n"; });
throw new Exception('callback-ok');
