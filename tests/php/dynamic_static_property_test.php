<?php

class DynamicStaticProperty_Target
{
    public static ?int $precision = 0;
}

class DynamicStaticProperty_Factory
{
    public function create(): DynamicStaticProperty_Target
    {
        return new DynamicStaticProperty_Target();
    }
}

$factory = new DynamicStaticProperty_Factory();

if ($factory->create()::$precision !== 0) {
    echo "FAIL dynamic static property\n";
    exit(1);
}

echo "PASS\n";
