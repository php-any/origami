<?php

class StaticPropSmoke
{
    protected static $resolver;

    public static function setResolver($v): void
    {
        static::$resolver = $v;
    }

    public static function getResolver()
    {
        return static::$resolver;
    }
}

StaticPropSmoke::setResolver('ok');
if (StaticPropSmoke::getResolver() !== 'ok') {
    echo "FAIL\n";
    exit(1);
}
echo "PASS\n";
