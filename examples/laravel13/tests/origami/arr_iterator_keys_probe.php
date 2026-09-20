<?php
require __DIR__.'/../../vendor/autoload.php';
$src = ['x-cloak' => true, 'class' => 'fi-topbar-open-sidebar-btn'];
$it = new ArrayIterator($src);
echo 'copy='.json_encode($it->getArrayCopy())."\n";
echo 'ita='.json_encode(iterator_to_array($it))."\n";
foreach ($it as $k => $v) {
    echo "fe type=".gettype($k)." k=".json_encode($k)." v=".json_encode($v)."\n";
}
$coll = new Illuminate\Support\Collection($src);
echo 'items='.json_encode($coll->all())."\n";
$it2 = $coll->getIterator();
echo 'coll_ita='.json_encode(iterator_to_array($it2))."\n";
