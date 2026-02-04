<?php

namespace tests\php;

use Generator;

// 1. array 一定是 iterable
Log::info('is_iterable([]): ' . (is_iterable([]) ? 'true' : 'false'));

// 2. stdClass 默认不可迭代（foreach 可以用，但 is_iterable 应为 false）
$obj = new \stdClass();
Log::info('is_iterable(stdClass): ' . (is_iterable($obj) ? 'true' : 'false'));

// 3. 自定义实现 Iterator 接口的对象应该是 iterable
class MyIterator implements \Iterator
{
    private array $data = [1, 2, 3];
    private int $position = 0;

    public function current(): mixed { return $this->data[$this->position]; }
    public function key(): mixed { return $this->position; }
    public function next(): void { ++$this->position; }
    public function rewind(): void { $this->position = 0; }
    public function valid(): bool { return isset($this->data[$this->position]); }
}

$it = new MyIterator();
Log::info('is_iterable(MyIterator): ' . (is_iterable($it) ? 'true' : 'false'));

// 4. 生成器（yield）也应视为 iterable
function gen(): Generator {
    yield 1;
    yield 2;
}

$g = gen();
Log::info('is_iterable(generator): ' . (is_iterable($g) ? 'true' : 'false'));

