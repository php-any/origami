<?php
use Illuminate\Support\Arr;

function check_arr_add($condition, $message) {
    if (!$condition) {
        throw new Exception($message);
    }
}

$config = ['driver' => 'mysql'];
$result = Arr::add(Arr::add($config, 'prefix', ''), 'name', 'default');
check_arr_add($result === ['driver' => 'mysql', 'prefix' => '', 'name' => 'default'], 'nested add must accept a temporary');
check_arr_add($config === ['driver' => 'mysql'], 'add must leave its input unchanged');

$source = ['value' => null, 'nested' => ['value' => null], 'other' => 7];
$alias = $source;
$result = Arr::add($source, 'value', 3);
check_arr_add($result['value'] === 3 && $source['value'] === null, 'existing null must be replaced only in the result');
$result = Arr::add($source, 'nested.value', 5);
check_arr_add($result['nested']['value'] === 5, 'nested null must be replaced');
check_arr_add($source === $alias && $source['nested']['value'] === null, 'nested updates must preserve the input and its alias');
$result = Arr::add($source, 'nested.new', 8);
check_arr_add($result['nested']['new'] === 8 && !array_key_exists('new', $source['nested']), 'nested insertion must isolate arrays');

foreach ([false, 0, '', []] as $existing) {
    check_arr_add(Arr::add(['value' => $existing], 'value', 9)['value'] === $existing, 'non-null values must be retained');
}
check_arr_add(Arr::add(['nested.value' => 4, 'nested' => []], 'nested.value', 9) === ['nested.value' => 4, 'nested' => []], 'literal dotted key takes precedence for lookup');
check_arr_add(Arr::add([], '', 2) === ['' => 2], 'empty string is a valid array key');
check_arr_add(Arr::add([], null, 2) === [], 'null key returns an existing non-null array');
check_arr_add(Arr::add(null, null, 2) === 2, 'null array and key return the supplied value');
check_arr_add(Arr::add(null, 'a.b', 2) === ['a' => ['b' => 2]], 'null input can become an array');
check_arr_add(Arr::add([], '0.value', 2) === [0 => ['value' => 2]], 'numeric path segments preserve integer keys');
check_arr_add(Arr::add(value: 2, key: 'value', array: []) === ['value' => 2], 'named arguments retain value semantics');
check_arr_add(Arr::add(...[[], 'value', 2]) === ['value' => 2], 'unpacked arguments accept temporary arrays');

$referenced = null;
$source = ['value' => &$referenced];
$result = Arr::add($source, 'value', 6);
check_arr_add($referenced === 6 && $result['value'] === 6, 'explicit PHP array references remain shared');

$source = ['value' => null];
$result = Arr::set($source, 'value', 10);
check_arr_add($source['value'] === 10 && $result['value'] === 10, 'set must remain a reference operation');
echo "Laravel Arr::add value semantics: PASS\n";
