<?php

use Illuminate\Support\Arr;
use Illuminate\Support\Collection;

foreach (['resources/css/filament/admin/theme.css', '', 42, 1.5, false, true] as $value) {
    if ((new Collection($value))->all() !== [$value] || collect($value)->all() !== [$value]) {
        throw new Exception('Collection must wrap a scalar without changing its value');
    }
    try {
        Arr::from($value);
        throw new Exception('Arr::from must still reject scalar values');
    } catch (InvalidArgumentException $expected) {}
}
if ((new Collection(null))->all() !== [] || (new Collection(['css' => 'theme.css']))->all() !== ['css' => 'theme.css']) {
    throw new Exception('Collection must preserve null and keyed array semantics');
}
enum ThemeMode { case Light; }
enum ThemeColor: string { case Indigo = 'indigo'; }
foreach ([ThemeMode::Light, ThemeColor::Indigo] as $value) {
    if ((new Collection($value))->all() !== [$value]) {
        throw new Exception('Collection must wrap enum cases without converting them');
    }
}
echo "collection scalar entrypoint OK\n";
