<?php
require __DIR__.'/../../vendor/autoload.php';

$c = new Illuminate\Support\Collection([
    'Illuminate\\View\\ViewServiceProvider',
    'App\\Providers\\AppServiceProvider',
]);
$p = $c->partition(fn ($x) => str_starts_with($x, 'Illuminate\\'));
echo 'partition='.$p->count()."\n";
foreach ($p as $i => $group) {
    echo "group_$i class=".get_class($group)." count=".$group->count()."\n";
    echo "  is_collection=".($group instanceof Illuminate\Support\Collection ? 'y' : 'n')."\n";
}

$pkg = ['BladeUI\\Heroicons\\BladeHeroiconsServiceProvider'];
$p->splice(1, 0, [$pkg]);
echo 'after_splice='.$p->count()."\n";
foreach ($p as $i => $group) {
    $t = is_array($group) ? 'array' : get_class($group);
    echo "slot_$i type=$t count=".(is_array($group) ? count($group) : $group->count())."\n";
}

$collapsed = $p->collapse();
echo 'collapsed='.$collapsed->count()."\n";
echo 'has_view='.(in_array('Illuminate\\View\\ViewServiceProvider', $collapsed->toArray(), true) ? 'y' : 'n')."\n";
