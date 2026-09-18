<?php

namespace tests\php;

/**
 * 对齐 Illuminate\Container::findInContextualBindings：
 * ArrayAccess 类上 $this->contextual[$concrete]['$connection'] = ...
 * 再用 $this->contextual[end($this->buildStack)][$abstract] ?? null 读回。
 */

class ContextualBindingEnd_Need
{
    public $connection;

    public function __construct(string $connection, ?int $chunkSize = null)
    {
        $this->connection = $connection;
    }
}

class ContextualBindingEnd_Box implements \ArrayAccess
{
    protected $contextual = [];
    protected $buildStack = [];

    public function offsetExists(mixed $offset): bool
    {
        return false;
    }

    public function offsetGet(mixed $offset): mixed
    {
        return null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
    }

    public function offsetUnset(mixed $offset): void
    {
    }

    public function add($concrete, $abstract, $impl)
    {
        $this->contextual[$concrete][$abstract] = $impl;
    }

    public function build($concrete)
    {
        $this->buildStack[] = $concrete;
        try {
            $param = (new \ReflectionClass($concrete))->getConstructor()->getParameters()[0];
            $abstract = '$'.$param->getName();
            $last = end($this->buildStack);
            if ($last !== $concrete) {
                Log::fatal('end($this->buildStack) 应为 '.$concrete.', 实际 '.var_export($last, true));
            }
            $binding = $this->contextual[$last][$abstract] ?? null;
            if ($binding === null) {
                Log::fatal('contextual 未命中 abstract='.$abstract.' last='.var_export($last, true));
            }

            return is_callable($binding) ? $binding() : $binding;
        } finally {
            array_pop($this->buildStack);
        }
    }
}

$box = new ContextualBindingEnd_Box();
$box->add(ContextualBindingEnd_Need::class, '$connection', fn () => 'sqlite');
$got = $box->build(ContextualBindingEnd_Need::class);
if ($got !== 'sqlite') {
    Log::fatal('contextual give 结果错误: '.var_export($got, true));
}

Log::info('contextual_binding_end 测试通过');
