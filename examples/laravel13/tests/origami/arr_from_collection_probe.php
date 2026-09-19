<?php

/**
 * Arr::from(Collection) 应对齐 PHP：摊平为 all()，不能包成 [Collection]。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$inner = new Illuminate\Support\Collection([1, 2, 3]);
$from = Illuminate\Support\Arr::from($inner);
echo 'arr_from_type='.gettype($from).' count='.(is_countable($from) ? count($from) : -1).' json='.json_encode($from)."\n";

$outer = new Illuminate\Support\Collection($inner);
$all = $outer->all();
echo 'outer_all_type='.gettype($all).' count='.(is_countable($all) ? count($all) : -1).' json='.json_encode($all)."\n";
echo 'outer_all0='.gettype($all[0] ?? null)."\n";

$keep = (new Illuminate\Support\Collection([[1, 2, 3]]))->all();
echo 'keep_inner='.json_encode($keep)."\n";
