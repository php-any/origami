<?php

namespace tests\php;

/**
 * parent:: 方法调用不应沿继承链重复执行。
 */

class ParentModify_Base {
    public int $n = 0;
    public function bump(): void {
        $this->n += 1;
    }
}

class ParentModify_Child extends ParentModify_Base {
    public function bump(): void {
        parent::bump();
    }
}

class ParentModify_Grand extends ParentModify_Child {
    public function bump(): void {
        parent::bump();
    }
}

$g = new ParentModify_Grand();
$g->bump();
if ($g->n !== 1) {
    Log::fatal("parent:: 链调用次数错误: n={$g->n}, 期望 1");
}

Log::info('parent:: 单次上溯测试通过');
