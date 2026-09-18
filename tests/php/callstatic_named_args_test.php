<?php
namespace tests\php;

/**
 * 模拟 Laravel Facade::__callStatic：命名实参须保留在 $parameters 数组键中，
 * 以便 $instance->$method(...$parameters) 能按名绑定目标方法形参。
 */
class FacadeTarget {
    public function renderHook(string $name, $scopes = null, array $data = []): string {
        return "name=$name;scopes=" . json_encode($scopes) . ";data=" . json_encode($data);
    }
}

class FacadeSim {
    public static function __callStatic($method, $parameters) {
        $instance = new FacadeTarget();
        return $instance->$method(...$parameters);
    }
}

$result = FacadeSim::renderHook('panels::auth.login.form.before', scopes: ['admin']);
$expected = 'name=panels::auth.login.form.before;scopes=["admin"];data=[]';
if ($result !== $expected) {
    Log::fatal("callstatic named args 失败: got [$result]");
}

Log::info('callstatic_named_args 测试通过');
