<?php

class ImplodeStringable
{
    public function __construct(private string $value)
    {
    }

    public function __toString(): string
    {
        return $this->value;
    }
}

$result = implode('-', [new ImplodeStringable('a'), new ImplodeStringable('b')]);
if ($result !== 'a-b') {
    Log::fatal('implode 未调用对象 __toString: ' . $result);
}

Log::info('implode_object_tostring_test OK');
