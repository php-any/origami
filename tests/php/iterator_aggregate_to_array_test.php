<?php

class MyAgg implements IteratorAggregate
{
    public function getIterator(): Traversable
    {
        return new ArrayIterator(['a', 'b', 'c']);
    }
}

$arr = iterator_to_array(new MyAgg(), false);
echo implode(',', $arr), PHP_EOL;
