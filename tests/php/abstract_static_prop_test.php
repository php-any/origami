<?php

namespace TestNs;

abstract class AbsModel
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

AbsModel::setResolver('ok');
if (AbsModel::getResolver() !== 'ok') {
    echo "FAIL abs\n";
    var_export(AbsModel::getResolver());
    echo "\n";
    exit(1);
}

echo "PASS\n";
