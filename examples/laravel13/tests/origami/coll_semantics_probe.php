<?php
/**
 * Collection 比较语义（where / unique / whereIn / sortBy 等）origami 与 PHP 对齐性探针。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;

function show(string $label, callable $fn): void
{
    echo $label.' => ';
    try {
        echo json_encode($fn(), JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
    } catch (Throwable $e) {
        echo 'EX '.get_class($e).' '.$e->getMessage();
    }
    echo "\n";
}

$data = [
    ['id' => 1, 'price' => 9,     'name' => 'a', 'flag' => 1],
    ['id' => 2, 'price' => 10,    'name' => 'a', 'flag' => '1'],
    ['id' => 3, 'price' => 100,   'name' => 'b', 'flag' => true],
    ['id' => 4, 'price' => '9.5', 'name' => 'c', 'flag' => null],
];

show('where_gt_int', fn () => collect($data)->where('price', '>', 9)->pluck('id')->all());
show('where_gt_str', fn () => collect($data)->where('price', '>', '9')->pluck('id')->all());
show('where_strict_int', fn () => collect($data)->where('flag', '===', 1)->pluck('id')->all());
show('where_strict_str', fn () => collect($data)->where('flag', '===', '1')->pluck('id')->all());
show('where_loose_bool', fn () => collect($data)->where('flag', true)->pluck('id')->all());
show('where_lt_100', fn () => collect($data)->where('price', '<', 100)->pluck('id')->all());
show('whereBetween', fn () => collect($data)->whereBetween('price', [10, 100])->pluck('id')->all());
show('whereNotBetween', fn () => collect($data)->whereNotBetween('price', [10, 100])->pluck('id')->all());
show('unique_name', fn () => collect($data)->unique('name')->pluck('id')->all());
show('unique_all', fn () => collect([1, '1', 1.0])->unique()->all());
show('uniqueStrict', fn () => collect([1, '1', 1.0])->uniqueStrict()->all());
show('uniqueStrict_name', fn () => collect($data)->uniqueStrict('name')->pluck('id')->all());
show('whereIn_loose', fn () => collect($data)->whereIn('flag', [1])->pluck('id')->all());
show('whereIn_strict', fn () => collect($data)->whereIn('flag', [1], true)->pluck('id')->all());
show('whereInStrict', fn () => collect($data)->whereInStrict('flag', [1])->pluck('id')->all());
show('whereNotIn', fn () => collect($data)->whereNotIn('flag', [1])->pluck('id')->all());
show('contains_strict', fn () => collect($data)->containsStrict('flag', 1));
show('sortBy_price', fn () => collect($data)->sortBy('price')->pluck('id')->all());
show('groupBy_flag', fn () => collect($data)->groupBy('flag')->keys()->all());
show('value_name', fn () => collect($data)->value('name'));
show('value_deep', fn () => collect([['a' => ['b' => 7]]])->value('a.b'));
show('toPrettyJson', fn () => collect([1, 2])->toPrettyJson());
show('forPage', fn () => collect([1, 2, 3, 4, 5])->forPage(2, 2)->all());
show('percentage', fn () => collect([1, 2, 3, 4])->percentage(fn ($v) => $v > 2));
show('pipeThrough', fn () => collect([1, 2, 3])->pipeThrough([fn ($c) => $c->map(fn ($v) => $v * 2), fn ($c) => $c->sum()]));
show('reduceInto', function () {
    $init = [];
    return collect([1, 2])->reduceInto($init, function (&$carry, $v, $k) { $carry[$k] = $v * 3; });
});
show('string_cast', fn () => 'c='.(string) collect([1, 2]));
show('hasMany_cb', fn () => collect([1, 2, 3])->hasMany(fn ($v) => $v > 1));
show('hasMany_false', fn () => collect([1, 2, 3])->hasMany(fn ($v) => $v > 2));
show('hasMany_args', fn () => collect($data)->hasMany('name', 'a'));
show('mapSpread', fn () => collect([[1, 2], [3, 4]])->mapSpread(fn ($a, $b, $k) => $a + $b + $k)->all());
show('eachSpread', function () {
    $seen = [];
    collect([[1, 2], [3, 4]])->eachSpread(function ($a, $b, $k) use (&$seen) { $seen[] = $a + $b + $k; });
    return $seen;
});
show('reduceInto_keys', function () {
    $init = [];
    return collect(['a' => 1, 'b' => 2])->reduceInto($init, function (&$carry, $v, $k) { $carry[$k] = $v * 3; });
});
show('reduceWithKeys', fn () => collect(['a' => 1, 'b' => 2])->reduceWithKeys(fn ($c, $v, $k) => $c.$k.$v, ''));
show('reduceSpread', fn () => collect([1, 2, 3])->reduceSpread(fn ($sum, $prod, $item) => [$sum + $item, $prod * $item], 0, 1));
show('mapToDictionary', fn () => collect($data)->mapToDictionary(fn ($i) => [$i['name'] => $i['id']])->map(fn ($g) => $g)->all());
show('mapToGroups', fn () => collect($data)->mapToGroups(fn ($i) => [$i['name'] => $i['id']])->map(fn ($g) => $g->all())->all());
show('pipeInto', fn () => collect([1, 2])->pipeInto(Illuminate\Support\Collection::class)->all());
show('fromJson', fn () => Collection::fromJson('{"a":1,"b":[2,3]}')->all());
show('escape_toJson', fn () => collect([1, 2])->escapeWhenCastingToString()->toJson());
show('no_escape_toJson', fn () => collect([1, 2])->escapeWhenCastingToString(false)->toJson());
show('percentage_empty', fn () => collect([])->percentage(fn ($v) => true));
show('unique_key_missing', fn () => collect([['a' => 1], ['b' => 2], ['a' => 1]])->unique('a')->all());
show('where_null_key', fn () => collect($data)->where('nope', '!=', 5)->pluck('id')->all());
show('where_explicit_null', fn () => collect($data)->where('flag', '=', null)->pluck('id')->all());
show('where_not_loose', fn () => collect($data)->where('flag', '!=', 1)->pluck('id')->all());
show('whereLeft_alias', fn () => collect($data)->where('name', 'a')->pluck('id')->all());
show('contains_missing_key', fn () => collect($data)->contains('nope', '!=', 5));
show('sort_regular', fn () => collect([9, 10, '9.5', 100])->sort()->values()->all());
show('sortAsc_helper', fn () => collect(['b' => 9, 'a' => 10])->sort()->all());
show('diff_loose', fn () => collect([1, 2, 3])->diff(['2'])->values()->all());
show('unique_after_diff', fn () => collect([1, 1, '1', 2])->uniqueStrict()->values()->all());
show('sum_float', fn () => collect([1, 2.5])->sum());
show('nth_keys', fn () => collect([1, 2, 3, 4])->nth(2)->all());
echo "DONE\n";
