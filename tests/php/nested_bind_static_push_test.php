<?php

namespace tests\php;

/**
 * 在 Closure::bind 上下文中触发「类方法注册的监听器」二次 push static 栈
 * （模拟 Login 视图 bind 内再 mount Notifications）。
 */

class NestedPush_Bus
{
    public array $listeners = [];

    public function on(callable $cb): void
    {
        $this->listeners[] = $cb;
    }

    public function trigger($target): void
    {
        foreach ($this->listeners as $cb) {
            $cb($target);
        }
    }
}

class NestedPush_Extend
{
    protected static $stack = [];

    public function boot(NestedPush_Bus $bus): void
    {
        $bus->on(function ($target) {
            $this->push($target);
        });
    }

    public function push($target): void
    {
        static::$stack[] = $target;
    }

    public static function current()
    {
        return end(static::$stack);
    }

    public static function reset(): void
    {
        static::$stack = [];
    }
}

class NestedPush_Login
{
    public string $name = 'Login';
}

class NestedPush_Notifications
{
    public string $name = 'Notifications';
}

NestedPush_Extend::reset();
$bus = new NestedPush_Bus();
(new NestedPush_Extend())->boot($bus);

$login = new NestedPush_Login();
$notif = new NestedPush_Notifications();

$bus->trigger($login);
if (NestedPush_Extend::current()->name !== 'Login') {
    Log::fatal('第一次 push 失败');
}

// 在 bind(Login) 内再次 trigger(Notifications)
$bound = \Closure::bind(function () use ($bus, $notif) {
    $bus->trigger($notif);
}, $login, $login);
$bound();

$cur = NestedPush_Extend::current();
if (!is_object($cur) || $cur->name !== 'Notifications') {
    Log::fatal('嵌套 bind 内二次 push 失败, current=' . (is_object($cur) ? $cur->name : var_export($cur, true)));
}

Log::info('嵌套 bind 内 static push 测试通过');
