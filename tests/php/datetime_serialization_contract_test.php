<?php
function dateWireCheck($ok, $message) { if (!$ok) throw new Exception($message); }
foreach ([DateTime::class, DateTimeImmutable::class] as $class) {
    $date = new $class('2026-01-02 03:04:05.123456', new DateTimeZone('UTC'));
    $copy = unserialize(serialize($date));
    dateWireCheck($copy instanceof $class && $copy !== $date, 'date identity');
    dateWireCheck($copy->format('Y-m-d H:i:s.u') === '2026-01-02 03:04:05.123456', 'date and microseconds');
    dateWireCheck($copy->__serialize() === ['date'=>'2026-01-02 03:04:05.123456', 'timezone_type'=>3, 'timezone'=>'UTC'], 'native date payload');
    dateWireCheck(serialize($copy) === serialize($date), 'date wire round trip');
}
echo "DateTime serialization OK\n";
