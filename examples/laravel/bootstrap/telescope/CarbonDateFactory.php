<?php

namespace Bootstrap\Telescope;

/**
 * DateFactory 适配：避免 Illuminate\Support\Carbon 的 parent::format 死循环。
 */
class CarbonDateFactory
{
    public function now($tz = null)
    {
        return \Carbon\Carbon::now($tz);
    }

    public function __call($method, $parameters)
    {
        return \Carbon\Carbon::$method(...$parameters);
    }
}
