<?php
$utc = new DateTimeImmutable('@0');
$local = $utc->setTimezone(new DateTimeZone('Asia/Shanghai'));
if ($local === $utc || $local->getTimestamp() !== 0 || $local->format('c') !== '1970-01-01T08:00:00+08:00' || $utc->format('c') !== '1970-01-01T00:00:00+00:00') { throw new RuntimeException('Timezone conversion must preserve instant and immutable source'); }
$zone = new DateTimeZone('America/New_York');
$a = (new DateTimeImmutable('@1730611800'))->setTimezone($zone);
$b = (new DateTimeImmutable('@1730615400'))->setTimezone($zone);
if ($a->format('H:i P') !== '01:30 -04:00' || $b->format('H:i P') !== '01:30 -05:00') { throw new RuntimeException('DST fold must preserve both instants'); }
echo "timezone OK\n";

class OverrideTimezone extends DateTimeZone {
    public function getName(): string { throw new RuntimeException('getName override called'); }
}
class CloneAwareDate extends DateTimeImmutable {
    public $copies = 0;
    public function __clone() { $this->copies++; }
}
$original = new CloneAwareDate('@0');
$converted = $original->setTimezone(new OverrideTimezone('Asia/Shanghai'));
if (!($converted instanceof CloneAwareDate) || $converted->copies !== 1 || $original->copies !== 0 || $converted->format('c') !== '1970-01-01T08:00:00+08:00') {
    throw new RuntimeException('Timezone conversion lost subclass or clone semantics');
}
$caught = false;
try { $original->setTimezone('UTC'); } catch (TypeError $error) { $caught = true; }
if (!$caught) { throw new RuntimeException('Invalid timezone argument accepted'); }
class ThrowingCloneDate extends DateTimeImmutable {
    public function __clone() { throw new RuntimeException('clone failure'); }
}
$caught = false;
try { (new ThrowingCloneDate('@0'))->setTimezone(new DateTimeZone('UTC')); }
catch (RuntimeException $error) { $caught = $error->getMessage() === 'clone failure'; }
if (!$caught) { throw new RuntimeException('Clone exception swallowed'); }
echo "timezone subclass OK\n";
