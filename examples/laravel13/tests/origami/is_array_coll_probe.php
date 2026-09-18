<?php
/**
 * is_array on Collection / ObjectValue semantics.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;
use Illuminate\Database\Eloquent\Collection as EC;
use App\Models\Order;

$c = new Collection([1,2]);
echo "support is_array=". (is_array($c) ? 'yes' : 'no') ." type=".gettype($c)."\n";
echo "support all is_array=". (is_array($c->all()) ? 'yes' : 'no') ."\n";

$e = Order::query()->limit(1)->get();
echo "eloquent is_array=". (is_array($e) ? 'yes' : 'no') ." class=".get_class($e)."\n";
$all = $e->all();
echo "eloquent all is_array=". (is_array($all) ? 'yes' : 'no') ." type=".gettype($all)."\n";
echo "eloquent all is_object=". (is_object($all) ? 'yes' : 'no') ."\n";
if (is_object($all)) {
    echo "all_class=".get_class($all)."\n";
}

// Arr::from
$x = Illuminate\Support\Arr::from($e);
echo "arr_from eloq is_array=". (is_array($x) ? 'yes' : 'no') ." type=".gettype($x)."\n";

// Nested construct
$e2 = new EC($e);
echo "nested construct ok count=".count($e2->all())."\n";

echo "DONE\n";
