<?php
function error_state_check($condition, $message) { if (!$condition) throw new Exception($message); }
$old = error_reporting(0);
error_state_check(error_reporting() === 0, 'reporting set/get');
error_clear_last();
$events = [];
$first = function ($level, $message) use (&$events) { $events[] = [$level, $message]; return true; };
error_state_check(set_error_handler($first, E_USER_WARNING | E_WARNING) === null, 'first handler');
error_state_check(trigger_error('warning', E_USER_WARNING) === true, 'warning return');
compact('error_state_missing');
$empty = []; $missing = $empty['error_state_key'];
error_state_check($events === [[E_USER_WARNING, 'warning'], [E_WARNING, 'compact(): Undefined variable $error_state_missing'], [E_WARNING, 'Undefined array key "error_state_key"']], 'warning dispatch');
error_state_check(error_get_last() === null, 'handled warnings do not set last error');
trigger_error('filtered', E_USER_NOTICE);
error_state_check(error_get_last()['message'] === 'filtered' && error_get_last()['type'] === E_USER_NOTICE, 'handler mask');
$second = function ($level, $message) { throw new RuntimeException($message); };
error_state_check(set_error_handler($second) === $first, 'previous identity');
try { trigger_error('from-handler', E_USER_NOTICE); throw new Exception('handler swallowed'); } catch (RuntimeException $e) { error_state_check($e->getMessage() === 'from-handler', 'exception propagation'); }
restore_error_handler();
trigger_error('restored', E_USER_WARNING);
error_state_check(count($events) === 4, 'restored handler mask');
set_error_handler(function () { return false; });
trigger_error('fallback', E_USER_WARNING);
error_state_check(error_get_last()['message'] === 'fallback', 'false falls back');
restore_error_handler();
restore_error_handler();
error_clear_last();
error_state_check(error_get_last() === null, 'clear');
try { trigger_error('bad', E_WARNING); throw new Exception('invalid level accepted'); } catch (ValueError $e) {}
error_reporting($old);
echo "error state contract OK\n";
