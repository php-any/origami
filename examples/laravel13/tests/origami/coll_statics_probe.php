<?php
/**
 * Collection 静态方法（times/range/unwrap）与 Gate 策略猜测路径。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Gate;

function show(string $label, callable $fn): void
{
    echo $label.' => ';
    try {
        echo json_encode($fn(), JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
    } catch (Throwable $e) {
        echo 'EX '.$e->getMessage();
    }
    echo "\n";
}

function policyName($policy)
{
    return is_object($policy) ? get_class($policy) : $policy;
}

show('times_null', fn () => Collection::times(3)->all());
show('times_cb', fn () => Collection::times(3, fn ($i) => $i * 10)->all());
show('times_cb_key', fn () => Collection::times(3, fn ($i, $k) => "$k:$i")->all());
show('times_zero', fn () => Collection::times(0)->all());
show('times_neg', fn () => Collection::times(-2)->all());
show('times_when_reverse', fn () => Collection::times(3)->when(true, fn ($c) => $c->concat([9]))->reverse()->values()->all());
show('times_chain_guard', fn () => Collection::times(3)->when(false, fn ($c) => $c->concat([9]))->reverse()->values()->all());
show('range', fn () => Collection::range(1, 5)->all());
show('unwrap_coll', fn () => Collection::unwrap(collect(['a' => 1])));
show('unwrap_arr', fn () => Collection::unwrap([1, 2]));
show('unwrap_null', fn () => Collection::unwrap(null));
show('unwrap_scalar', fn () => Collection::unwrap(5));
show('gate_policy_for_order_item', fn () => policyName(Gate::getPolicyFor(App\Models\OrderItem::class)));
show('gate_policy_for_user', fn () => policyName(Gate::getPolicyFor(App\Models\User::class)));
show('gate_policy_missing', fn () => Gate::getPolicyFor('App\\Models\\Nowhere\\Ghost'));
show('str_password_len', fn () => strlen(Illuminate\Support\Str::password(16)));
echo "DONE\n";
