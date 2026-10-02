<?php

use Illuminate\Support\Arr;

$profile = (object) ['name' => 'profile', 'sort' => -1];
$logout = (object) ['name' => 'logout', 'sort' => 1];
$items = Arr::collapse([['profile' => $profile], ['logout' => $logout]]);
$groups = collect($items)->groupBy(fn ($item) => $item->sort < 0, preserveKeys: true)->all();
if (!$groups[true]->has('profile') || Arr::first($groups[true]) !== $profile || Arr::last($groups[false]) !== $logout) {
    throw new Exception('Filament user menu must retain action keys and first action identity');
}
if (Arr::first(collect(['x' => 1, 'y' => 2]), fn ($value, $key) => $key === 'y') !== 2
    || Arr::first(collect(), null, fn () => 'empty') !== 'empty') {
    throw new Exception('Arr::first must support collection callbacks and lazy defaults');
}
echo "collection first user menu OK\n";
