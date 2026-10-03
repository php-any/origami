<?php
function callable_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
class CallbackScope {
    private function hidden() { return 'private'; }
    public function instance() { return 'instance'; }
    public static function visible() { return 'static'; }
    public function inside() { return is_callable([$this, 'hidden']); }
}
$scope = new CallbackScope;
callable_assert(!is_callable([$scope, 'hidden']) && $scope->inside(), 'private visibility');
callable_assert(!is_callable([CallbackScope::class, 'instance']), 'static instance validation');
$name = null;
callable_assert(is_callable('CallbackScope::visible', false, $name) && $name === 'CallbackScope::visible', 'callable string name');
callable_assert(is_callable(['MissingCallbackClass', 'method'], true, $name) && $name === 'MissingCallbackClass::method', 'syntax-only should not load');
callable_assert(!is_callable(['MissingCallbackClass', 'method']), 'missing class callable');
class MagicCallbackScope {
    private function hidden() { return 'hidden'; }
    public function __call($name, $arguments) { return [$name, $arguments]; }
    public static function __callStatic($name, $arguments) { return [$name, $arguments]; }
}
$magic = new MagicCallbackScope;
callable_assert(is_callable([$magic, 'missing']) && is_callable([$magic, 'hidden']), 'instance magic callable');
callable_assert(is_callable('MagicCallbackScope::missing'), 'static magic callable');
callable_assert(call_user_func([$magic, 'missing'], 1, 2) === ['missing', [1, 2]], 'instance magic invocation');
callable_assert(call_user_func('MagicCallbackScope::missing', 3) === ['missing', [3]], 'static magic invocation');
callable_assert(call_user_func_array([$magic, 'missing'], ['named' => 7]) === ['missing', ['named' => 7]], 'magic named arguments');
function accepts_callback(callable $callback) { return true; }
callable_assert(accepts_callback([$magic, 'missing']), 'callable type magic');
$caught = false;
try { accepts_callback([$scope, 'hidden']); } catch (TypeError $error) { $caught = true; }
callable_assert($caught, 'callable type visibility');
echo "callable visibility magic ok\n";
