<?php
$source = new DateTime('@1234567890');
$copy = DateTimeImmutable::createFromInterface($source);
if ($copy === $source || get_class($copy) !== DateTimeImmutable::class || $copy->getTimestamp() !== 1234567890) { throw new Exception('DateTime interface copy'); }
$source->setTimestamp(42);
if ($copy->getTimestamp() !== 1234567890) { throw new Exception('DateTime snapshot'); }
$mutable = DateTime::createFromInterface($copy);
if (get_class($mutable) !== DateTime::class || $mutable === $copy || $mutable->getTimestamp() !== 1234567890) { throw new Exception('DateTime mutable interface copy'); }
echo "DateTime interface copies OK\n";
