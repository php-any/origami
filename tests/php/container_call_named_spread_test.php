<?php
namespace tests\php;

/**
 * 模拟 Laravel BoundMethod::call 对 Closure 的调用模式：
 * $callback(...array_values($dependencies))
 * 其中 dependencies 由关联数组按参数名解析后合并为位置数组。
 */
function container_call_simulate(callable $callback, array $parameters): mixed
{
    // 简化：已知参数名 data/scopes，按反射顺序组装（与 getMethodDependencies 等价路径）
    $deps = [];
    if (array_key_exists('data', $parameters)) {
        $deps[] = $parameters['data'];
        unset($parameters['data']);
    }
    if (array_key_exists('scopes', $parameters)) {
        $deps[] = $parameters['scopes'];
        unset($parameters['scopes']);
    }
    $deps = array_merge($deps, array_values($parameters));
    return $callback(...array_values($deps));
}

$hook = function ($data, $scopes = null) {
    return "data=" . json_encode($data) . ";scopes=" . json_encode($scopes);
};

$result = container_call_simulate($hook, ['data' => ['x' => 1], 'scopes' => ['a', 'b']]);
if ($result !== 'data={"x":1};scopes=["a","b"]') {
    Log::fatal('container_call spread 失败: ' . $result);
}

// 直接 spread 关联数组（错误路径：应把键当命名参数，而非仅取值）
$assoc = ['data' => ['k' => 2], 'scopes' => null];
try {
    $bad = $hook(...$assoc);
    // PHP 会把关联键当命名参数
    if ($bad !== 'data={"k":2};scopes=null') {
        Log::fatal('assoc spread 命名绑定失败: ' . $bad);
    }
} catch (Throwable $e) {
    Log::fatal('assoc spread 异常: ' . $e->getMessage());
}

Log::info('container_call_named_spread 测试通过');
