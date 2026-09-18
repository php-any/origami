<?php
/**
 * 复现 Eloquent/Support Collection::all 深度问题。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;
use Illuminate\Database\Eloquent\Collection as EloquentCollection;
use App\Models\Order;

echo "support_empty\n";
$c = new Collection();
$a = $c->all();
echo "type=".gettype($a)." count=".(is_array($a)?count($a):-1)."\n";

echo "support_nested\n";
$c2 = new Collection(new Collection([1, 2, 3]));
$a2 = $c2->all();
echo "type=".gettype($a2)." count=".(is_array($a2)?count($a2):-1)." is_coll=".(is_object($a2)?get_class($a2):'-')."\n";

echo "eloquent_from_query\n";
try {
    $orders = Order::query()->latest()->limit(5)->get();
    echo "get_class=".get_class($orders)."\n";
    $all = $orders->all();
    echo "all_type=".gettype($all)." count=".(is_countable($all)?count($all):-1)."\n";
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n";
}

echo "arr_from\n";
try {
    $x = Illuminate\Support\Arr::from(new Collection(['a' => 1]));
    echo "arr_from_type=".gettype($x)."\n";
} catch (Throwable $e) {
    echo "arr_from EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
