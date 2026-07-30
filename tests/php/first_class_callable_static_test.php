<?php

namespace tests\php;

/**
 * PHP 8.1 first-class callable $this->method(...) 必须绑定 $this，
 * 使回调内 static:: 可用（Laravel AliasLoader::register 依赖此语义）。
 */

class FirstClass_AliasLoader
{
    protected $aliases;
    protected $registered = false;
    protected static $facadeNamespace = 'Facades\\';
    protected static $instance;

    private function __construct($aliases)
    {
        $this->aliases = $aliases;
    }

    public static function getInstance(array $aliases = [])
    {
        if (is_null(static::$instance)) {
            return static::$instance = new static($aliases);
        }

        return static::$instance;
    }

    public function load($alias)
    {
        if (static::$facadeNamespace && str_starts_with($alias, static::$facadeNamespace)) {
            return 'facade';
        }

        if (isset($this->aliases[$alias])) {
            return 'alias';
        }

        return null;
    }

    public function register()
    {
        if (! $this->registered) {
            spl_autoload_register($this->load(...), true, true);
            $this->registered = true;
        }
    }
}

$loader = FirstClass_AliasLoader::getInstance(['App' => 'X']);
$cb = $loader->load(...);
if (!is_callable($cb)) {
    Log::fatal('first-class callable 未返回可调用');
}

$r = $cb('App');
if ($r !== 'alias') {
    Log::fatal('直接调用 first-class 失败: '.var_export($r, true));
}

$r2 = $cb('Facades\\Foo');
if ($r2 !== 'facade') {
    Log::fatal('static::$facadeNamespace 在 first-class 中失败: '.var_export($r2, true));
}

$loader->register();
class_exists('Some\\MissingFirstClassXYZ', true);

Log::info('first_class_callable_static 测试通过');
