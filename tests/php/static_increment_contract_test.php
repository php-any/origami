<?php
function staticIncrementCheck($ok) { if (!$ok) throw new Exception('static increment contract'); }
class StaticCounter {
    public static int $count = 4;
    public static function update() { return [self::$count++, ++self::$count, self::$count--, --self::$count]; }
}
staticIncrementCheck(StaticCounter::update() === [4,6,6,4]);
$alias =& StaticCounter::$count;
staticIncrementCheck(StaticCounter::$count++ === 4 && $alias === 5);
staticIncrementCheck(--StaticCounter::$count === 4 && $alias === 4);
class StaticDerived extends StaticCounter {
    public static int $count = 10;
    public static function late() { return static::$count++; }
}
staticIncrementCheck(StaticDerived::late() === 10 && StaticDerived::$count === 11 && StaticCounter::$count === 4);
echo "static increment contract OK\n";
