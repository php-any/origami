<?php

namespace tests\php;

/**
 * new self 必须词法绑定到方法定义类（非 late static），对齐 PHP self vs static。
 * 复现 Laravel Collection::toBase / Eloquent::pluck 无限递归根因。
 */
class NewSelfLex_Base
{
    public $payload;

    public function __construct($payload = null)
    {
        $this->payload = $payload;
    }

    public function makeSelf()
    {
        return new self();
    }

    public function makeStatic()
    {
        return new static();
    }

    public function toBaseLike()
    {
        return new self($this);
    }

    public function tag()
    {
        return 'base';
    }
}

class NewSelfLex_Child extends NewSelfLex_Base
{
    public function pluckLike()
    {
        // 若 new self 误成 Child，会再次进 pluckLike → 无限递归
        return $this->toBaseLike()->tag();
    }

    public function tag()
    {
        return 'child';
    }
}

$child = new NewSelfLex_Child('x');
$selfInst = $child->makeSelf();
if (!($selfInst instanceof NewSelfLex_Base) || $selfInst instanceof NewSelfLex_Child) {
    \Log::fatal('new self 应为 Base，实际: '.get_class($selfInst));
}

$staticInst = $child->makeStatic();
if (!($staticInst instanceof NewSelfLex_Child)) {
    \Log::fatal('new static 应为 Child，实际: '.get_class($staticInst));
}

$baseFromChild = $child->toBaseLike();
if (get_class($baseFromChild) !== NewSelfLex_Base::class) {
    \Log::fatal('toBaseLike new self 应为 Base，实际: '.get_class($baseFromChild));
}

$tag = $child->pluckLike();
if ($tag !== 'base') {
    \Log::fatal('pluckLike 经 toBaseLike 后 tag 应为 base，实际: '.$tag);
}

\Log::info('new_self_lexical_binding 测试通过');
