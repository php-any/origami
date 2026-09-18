<?php

namespace tests\syntax;

/**
 * PHP 8.0：构造器属性提升。
 */

class SyntaxPromo_User
{
    public function __construct(
        public string $name,
        private int $id = 1,
    ) {
    }

    public function id(): int
    {
        return $this->id;
    }
}

$u = new SyntaxPromo_User('alice');
if ($u->name !== 'alice') {
    Log::fatal('提升 public 属性失败');
}
if ($u->id() !== 1) {
    Log::fatal('提升 private 默认值失败');
}

Log::info('syntax constructor promotion 测试通过');
