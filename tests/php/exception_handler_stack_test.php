<?php
function handlerStackCheck($ok, $message) {
    if (!$ok) throw new Exception($message);
}
$first = function (Throwable $e) {};
$second = function (Throwable $e) {};
set_exception_handler($first, ...[]);
restore_exception_handler();
handlerStackCheck(restore_exception_handler() === true, 'empty restore');
handlerStackCheck(set_exception_handler($first) === null, 'first registration');
handlerStackCheck(set_exception_handler($second) === $first, 'previous identity');
handlerStackCheck(set_exception_handler(null) === $second, 'disable');
handlerStackCheck(restore_exception_handler() === true, 'restore second');
handlerStackCheck(set_exception_handler($first) === $second, 'second restored');
restore_exception_handler();
restore_exception_handler();
handlerStackCheck(set_exception_handler(null) === $first, 'first restored');
restore_exception_handler();
restore_exception_handler();
handlerStackCheck(set_exception_handler(null) === null, 'default restored');
restore_exception_handler();
foreach ([123, 'missing_handler_function', ['MissingHandlerClass', 'method'], [2 => 'X', 3 => 'method']] as $bad) {
    try {
        set_exception_handler($bad);
        throw new Exception('accepted invalid callback');
    } catch (TypeError $e) {}
}
try {
    set_exception_handler();
    throw new Exception('accepted missing callback');
} catch (ArgumentCountError $e) {}
try {
    set_exception_handler($first, 1);
    throw new Exception('accepted extra callback argument');
} catch (ArgumentCountError $e) {}
echo "handler-stack-ok\n";
