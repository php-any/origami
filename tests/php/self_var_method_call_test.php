<?php

namespace tests\php;

/**
 * PHP：self::$method() / static::$method() 用变量值作为静态方法名。
 * Symfony IpUtils::checkIp 依赖 self::$method($ip, $cidr)。
 */

class SelfVarMethod_Host
{
    public static function checkIp4(string $ip, string $mask): string
    {
        return "v4:$ip/$mask";
    }

    public static function checkIp6(string $ip, string $mask): string
    {
        return "v6:$ip/$mask";
    }

    public static function checkIp(string $requestIp, string $cidr): string
    {
        $method = substr_count($requestIp, ':') > 1 ? 'checkIp6' : 'checkIp4';
        return self::$method($requestIp, $cidr);
    }
}

$v4 = SelfVarMethod_Host::checkIp('127.0.0.1', '127.0.0.0/8');
if ($v4 !== 'v4:127.0.0.1/127.0.0.0/8') {
    Log::fatal('self::$method() IPv4 失败: ' . var_export($v4, true));
}

$v6 = SelfVarMethod_Host::checkIp('::1', '::1/128');
if ($v6 !== 'v6:::1/::1/128') {
    Log::fatal('self::$method() IPv6 失败: ' . var_export($v6, true));
}

Log::info('self::$method() 动态静态方法调用测试通过');
