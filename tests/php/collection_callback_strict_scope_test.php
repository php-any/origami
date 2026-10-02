<?php
declare(strict_types=1);

$callback = function (string $value, $key): string { return $key.':'.$value; };
// Collection::map forwards to Arr::map / array_map, which bind callback
// arguments in weak mode even when the closure is declared in a strict file.
$items = ['enabled' => true, 'disabled' => false, 'number' => 12];
$expected = ['enabled' => 'enabled:1', 'disabled' => 'disabled:', 'number' => 'number:12'];
if (collect($items)->map($callback)->all() !== $expected
    || Illuminate\Support\Arr::map($items, $callback) !== $expected) {
    throw new Exception('Collection callback inherited the strict caller mode');
}
if (collect([true])->mapWithKeys(fn (string $value): array => ['key' => $value])->all() !== ['key' => '1']) {
    throw new Exception('Arr PHP callback must use its weak declaration mode');
}
try {
    $callback(true, 'direct');
    throw new Exception('Direct call must preserve strict parameter checking');
} catch (TypeError $expectedError) {}
try {
    collect([true])->map(fn (string $value): string => true);
    throw new Exception('Callback return must preserve its strict declaration mode');
} catch (TypeError $expectedError) {}
echo "collection callback strict scope OK\n";
