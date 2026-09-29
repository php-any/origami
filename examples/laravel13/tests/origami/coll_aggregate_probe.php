<?php
/**
 * Collection 聚合方法（sum/avg/max/min/median/nth/sort 等）origami 与 PHP 对齐性探针。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

function t($l, $fn)
{
    echo $l.' => ';
    try {
        echo json_encode($fn());
    } catch (Throwable $e) {
        echo 'EX '.get_class($e).': '.$e->getMessage();
    }
    echo "\n";
}

t('sum_int', fn () => collect([1, 2, 3])->sum());
t('sum_float', fn () => collect([1, 2.5])->sum());
t('sum_numeric_str', fn () => collect(['1', '2.5'])->sum());
t('sum_empty_arr', fn () => collect([])->sum());
t('sum_null', fn () => collect([null, 1])->sum());
t('sum_empty_str', fn () => collect(['', 1])->sum());
t('sum_key', fn () => collect([['p' => '3'], ['p' => 4]])->sum('p'));
t('sum_cb', fn () => collect([1, 2])->sum(fn ($v) => $v * 10));
t('avg_int', fn () => collect([1, 2, 3])->avg());
t('avg_div', fn () => collect([1, 2])->avg());
t('avg_null', fn () => collect([null, 2, 4])->avg());
t('avg_empty', fn () => collect([])->avg());
t('max_mixed', fn () => collect([1, '9.5', 10])->max());
t('max_str', fn () => collect(['b', 'a'])->max());
t('max_null', fn () => collect([null, 3])->max());
t('max_empty', fn () => collect([])->max());
t('min_str', fn () => collect(['b', 'a'])->min());
t('min_key', fn () => collect([['p' => 2], ['p' => 1]])->min('p'));
t('median_even', fn () => collect([3, 1, 2, 4])->median());
t('median_odd', fn () => collect([3, 1, 2])->median());
t('median_str', fn () => collect(['3', '1'])->median());
t('median_empty', fn () => collect([])->median());
t('mode', fn () => collect([1, 1, 2])->mode());
t('mode_key', fn () => collect([['p' => 1], ['p' => 1], ['p' => 2]])->mode('p'));
t('nth_keys', fn () => collect([1, 2, 3, 4])->nth(2)->all());
t('nth_offset', fn () => collect([1, 2, 3, 4])->nth(2, 1)->all());
t('sort_numeric', fn () => collect([10, 9, 100])->sort()->values()->all());
t('countBy', fn () => collect([1, 1, 2])->countBy()->all());
t('sum_by_collection', fn () => collect([['price' => '1.5'], ['price' => '2.5']])->sum('price'));
echo "DONE\n";
