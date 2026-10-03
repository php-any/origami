<?php
class StaticArrayOwner {
    public static array $paths = [];
    protected static array $groups = [];
    public static function add($class, $group, $path) {
        if (!array_key_exists($class, static::$paths)) { static::$paths[$class] = []; }
        static::$paths[$class] = array_merge(static::$paths[$class], [$path]);
        if (!array_key_exists($group, static::$groups)) { static::$groups[$group] = []; }
        static::$groups[$group][] = $path;
    }
    public static function groups() { return static::$groups; }
}
class StaticArrayChild extends StaticArrayOwner {}
StaticArrayChild::add('provider', 'assets', 'one');
$snapshot = StaticArrayOwner::$paths;
StaticArrayChild::add('provider', 'assets', 'two');
if ($snapshot !== ['provider' => ['one']] || StaticArrayOwner::$paths !== ['provider' => ['one', 'two']]) { throw new Exception('static array write lost storage or COW'); }
if (StaticArrayChild::groups() !== ['assets' => ['one', 'two']]) { throw new Exception('inherited static nested write lost storage'); }
$calls = 0;
function static_array_owner_name() { global $calls; $calls++; return StaticArrayOwner::class; }
$copy = StaticArrayOwner::$paths;
(static_array_owner_name())::$paths['other'] = ['three'];
if ($calls !== 1 || isset($copy['other']) || StaticArrayOwner::$paths['other'] !== ['three']) { throw new Exception('dynamic static receiver evaluation or COW'); }
unset(StaticArrayOwner::$paths['other']);
if (isset(StaticArrayOwner::$paths['other'])) { throw new Exception('static unset lost owner'); }
echo "static array dimension owners OK\n";
