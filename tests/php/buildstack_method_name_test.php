<?php

namespace tests\php;

/**
 * 方法名 build 与属性 buildStack：子类覆盖 build 后 $this->buildStack[] 仍应写属性。
 */
class BuildStackName_Base
{
    protected $buildStack = [];

    public function build($c)
    {
        $this->buildStack[] = $c;
        return $this->buildStack;
    }
}

class BuildStackName_Child extends BuildStackName_Base
{
    public function build($c)
    {
        $t = gettype($this->buildStack);
        if ($t !== 'array') {
            Log::fatal('子类 build() 内 $this->buildStack 应为 array, 实际 '.$t);
        }
        return parent::build($c);
    }
}

$c = new BuildStackName_Child();
$out = $c->build('x');
if (!is_array($out) || ($out[0] ?? null) !== 'x') {
    Log::fatal('parent::build 追加失败: '.var_export($out, true));
}

Log::info('buildstack_method_name 测试通过');
