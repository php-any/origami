<?php
namespace tests\php;

// 测试：Generator spread 展开到函数调用
function spread_gen_sum(...$nums) {
    $sum = 0;
    foreach ($nums as $n) {
        $sum += $n;
    }
    return $sum;
}

function spread_gen_yield() {
    yield 1;
    yield 2;
    yield 3;
}

$r1 = spread_gen_sum(...spread_gen_yield());
if ($r1 !== 6) Log::fatal("generator spread 到函数失败: $r1");

// 测试：Generator spread 展开到数组字面量
$arr = [0, ...spread_gen_yield(), 10];
if (count($arr) !== 5) Log::fatal("generator spread 到数组 count 失败");
if (array_sum($arr) !== 16) Log::fatal("generator spread 到数组 sum 失败");

// 测试：Generator spread 到构造函数（模拟 doctrine/inflector 场景）
class SpreadGen_Pattern {
    private $pattern;
    public function __construct($p) { $this->pattern = $p; }
    public function getPattern() { return $this->pattern; }
}

function spread_gen_patterns() {
    yield new SpreadGen_Pattern('.*ss');
    yield new SpreadGen_Pattern('data');
    yield new SpreadGen_Pattern('fuchsia');
}

class SpreadGen_Patterns {
    private $patterns;
    public function __construct(SpreadGen_Pattern ...$patterns) {
        $this->patterns = array_map(function(SpreadGen_Pattern $pattern) {
            return $pattern->getPattern();
        }, $patterns);
    }
    public function getPatterns() { return $this->patterns; }
}

$patterns = new SpreadGen_Patterns(...spread_gen_patterns());
$result = $patterns->getPatterns();
if (count($result) !== 3) Log::fatal("generator spread 到构造 count 失败");
if ($result[0] !== '.*ss' || $result[1] !== 'data' || $result[2] !== 'fuchsia') {
    Log::fatal("generator spread 到构造值失败");
}

// 测试：func_get_args 中 Generator spread
function spread_gen_fga(...$args) {
    return func_get_args();
}
$fga = spread_gen_fga(...spread_gen_yield());
if (count($fga) !== 3) Log::fatal("func_get_args generator spread 失败");

// 测试：yield from 的 Generator 展开
function spread_gen_from() {
    yield from spread_gen_yield();
    yield 4;
}
$r2 = spread_gen_sum(...spread_gen_from());
if ($r2 !== 10) Log::fatal("yield from generator spread 失败: $r2");

Log::info('spread_generator 测试通过');
