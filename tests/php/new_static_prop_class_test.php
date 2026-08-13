<?php

namespace tests\php;

/**
 * new static::$prop(...) 动态类名（Laravel Application::configure）。
 */

class NewStaticProp_Builder
{
    public $app;

    public function __construct($app)
    {
        $this->app = $app;
    }

    public function withKernels()
    {
        return $this;
    }
}

class NewStaticProp_App
{
    public static $applicationBuilder = NewStaticProp_Builder::class;

    public function __construct($base = null)
    {
    }

    public static function configure()
    {
        return (new static::$applicationBuilder(new static('base')))->withKernels();
    }
}

$x = NewStaticProp_App::configure();
if (!($x instanceof NewStaticProp_Builder)) {
    \Log::fatal('new static::$prop 应得到 Builder 实例');
}

\Log::info('new_static_prop_class 测试通过');
