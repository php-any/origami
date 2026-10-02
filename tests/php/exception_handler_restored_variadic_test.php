<?php
class RestoredVariadicHandler {
    public static function handle(Throwable ...$errors) {
        echo 'restored:' . count($errors) . ':' . $errors[0]->getMessage() . "\n";
    }
}
set_exception_handler([1 => 'handle', 0 => 'RestoredVariadicHandler']);
set_exception_handler(function () { echo "wrong-handler\n"; });
restore_exception_handler();
throw new Exception('variadic-ok');
