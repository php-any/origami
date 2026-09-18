<?php

namespace tests\syntax;

/**
 * PHP 8.0：nullsafe ?->。
 */

class SyntaxNs_Node
{
    public $child;
    public $name;

    public function __construct($child = null, $name = 'n')
    {
        $this->child = $child;
        $this->name = $name;
    }
}

$root = new SyntaxNs_Node();
if ($root?->child?->name !== null) {
    Log::fatal('nullsafe 链应对 null 返回 null');
}

$root = new SyntaxNs_Node(new SyntaxNs_Node(null, 'leaf'));
if ($root?->child?->name !== 'leaf') {
    Log::fatal('nullsafe 有值失败');
}

Log::info('syntax nullsafe 测试通过');
