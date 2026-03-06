<?php

namespace tests\php;

/**
 * 测试 trait 内的 use Trait1, Trait2 语法（trait 组合）
 * 用于验证 Carbon\Traits\Date 等使用子 trait 的场景
 */

trait TraitUse_Base {
    public function baseMethod(): string {
        return 'base';
    }
}

trait TraitUse_Extra {
    public function extraMethod(): string {
        return 'extra';
    }
}

trait TraitUse_Composite {
    use TraitUse_Base;
    use TraitUse_Extra;

    public function compositeMethod(): string {
        return $this->baseMethod() . '+' . $this->extraMethod();
    }
}

class TraitUse_TestClass {
    use TraitUse_Composite;
}

$obj = new TraitUse_TestClass();
$result = $obj->compositeMethod();
if ($result !== 'base+extra') {
    throw new \Exception('trait use 组合测试失败: 期望 base+extra, 实际 ' . $result);
}
\Log::info('trait use 组合测试通过');
