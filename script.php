<?php

// 测试 $target->{$segment} 的代码（使用类属性）

class Target {
    public string $name;
}

$segment = 'name';

$target = new Target();
$target->name = 'origami';

echo $target->{$segment};
