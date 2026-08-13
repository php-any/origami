<?php

namespace tests\php;

/**
 * 对象 === 同一性：静态缓存实例与方法内 $this 应判为同一对象。
 */

class ObjectIdentity_Factory
{
    private static ?self $defaultInstance = null;
    public $testNow = null;
    public $id;

    public function __construct()
    {
        static $n = 0;
        $this->id = ++$n;
    }

    public static function getDefaultInstance(): self
    {
        return self::$defaultInstance ??= new self();
    }

    public function getTestNow()
    {
        if ($this->testNow === null) {
            $factory = self::getDefaultInstance();
            if ($factory !== $this) {
                return $factory->getTestNow();
            }
        }
        return $this->testNow;
    }
}

$a = ObjectIdentity_Factory::getDefaultInstance();
$b = ObjectIdentity_Factory::getDefaultInstance();
if ($a !== $b) {
    \Log::fatal('两次 getDefaultInstance 应返回同一实例');
}

$r = $a->getTestNow();
if ($r !== null) {
    \Log::fatal('getTestNow 应返回 null，不应无限递归');
}

\Log::info('object_identity_strict 测试通过');
