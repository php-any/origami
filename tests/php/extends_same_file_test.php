<?php

class ExtendsSameFile_Parent
{
    public function hello(): string
    {
        return 'parent';
    }
}

class ExtendsSameFile_Child extends ExtendsSameFile_Parent
{
    public function greet(): string
    {
        return $this->hello();
    }
}

$c = new ExtendsSameFile_Child();
if ($c->greet() !== 'parent') {
    echo "FAIL greet\n";
    exit(1);
}

echo "PASS\n";
