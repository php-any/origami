<?php

namespace tests\php;

interface HasStaticProperty
{
    public static string $name = 'iface-prop';
}

// 读取接口静态属性，确认 node/call_static_property.go 能通过 GetOrLoadInterface 找到接口
Log::info('interface static property: ' . HasStaticProperty::$name);

