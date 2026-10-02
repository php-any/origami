<?php
declare(strict_types=1);
class CallbackOwner {
    public string $prefix = 'object:';
    public function mapped(int $number): string { return $this->prefix . $number; }
}
$owner = new CallbackOwner();
$callback = function (int $number): string { return $this->prefix . $number; };
$bound = $callback->bindTo($owner, CallbackOwner::class);
if (array_map($bound, ['2']) !== ['object:2']) { throw new Exception('bound callback context'); }
if (array_map([$owner, 'mapped'], ['3']) !== ['object:3']) { throw new Exception('method callback context'); }
$defaults = array_map(fn(int $number, string $suffix = 'default') => $number . $suffix, ['4']);
if ($defaults !== ['4default']) { throw new Exception('native callback default'); }
try {
    array_map(fn(int $number) => $number, [null]);
    throw new Exception('missing null TypeError');
} catch (TypeError $error) {}
try {
    array_map(fn(int $number, int $required) => $required, [1]);
    throw new Exception('missing ArgumentCountError');
} catch (ArgumentCountError $error) {}
echo "typed native callback: PASS\n";
