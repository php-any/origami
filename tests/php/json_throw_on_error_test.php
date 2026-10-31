<?php
function check_json_throw($condition, $message) {
    if (!$condition) { throw new RuntimeException($message); }
}

json_decode('{');
check_json_throw(json_last_error() === JSON_ERROR_SYNTAX, 'initial last error');
foreach ([['{', 512, JSON_ERROR_SYNTAX], ['[[1]]', 2, JSON_ERROR_DEPTH]] as $case) {
    $caught = false;
    try { json_decode($case[0], true, $case[1], JSON_THROW_ON_ERROR); }
    catch (JsonException $e) {
        $caught = true;
        check_json_throw($e instanceof Exception, 'exception inheritance');
        check_json_throw($e->getCode() === $case[2], 'decode error code');
        check_json_throw($e->getMessage() !== '', 'exception message');
    }
    check_json_throw($caught, 'decode must throw');
    check_json_throw(json_last_error() === JSON_ERROR_SYNTAX, 'throw preserves previous error');
}
check_json_throw(json_decode('{"x":1}', true, 512, JSON_THROW_ON_ERROR)['x'] === 1, 'valid decode');
check_json_throw(json_last_error() === JSON_ERROR_SYNTAX, 'successful throwing decode preserves error');
check_json_throw(json_encode(['ok' => true], JSON_THROW_ON_ERROR) === '{"ok":true}', 'valid encode');
check_json_throw(json_last_error() === JSON_ERROR_SYNTAX, 'successful throwing encode preserves error');
$recursive = [];
$recursive['self'] = &$recursive;
$caught = false;
try { json_encode($recursive, JSON_THROW_ON_ERROR); }
catch (JsonException $e) {
    $caught = true;
    check_json_throw($e->getCode() === JSON_ERROR_RECURSION, 'encode error code');
}
check_json_throw($caught, 'encode must throw');
check_json_throw(json_last_error() === JSON_ERROR_SYNTAX, 'encode throw preserves error');
json_decode('null');
check_json_throw(json_last_error() === JSON_ERROR_NONE, 'normal success clears error');
$manual = new JsonException('manual', 9);
check_json_throw($manual->getMessage() === 'manual' && $manual->getCode() === 9, 'manual construction');
echo "json_throw_on_error_test: OK\n";
