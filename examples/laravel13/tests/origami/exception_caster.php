<?php
require __DIR__.'/../../vendor/autoload.php';

$exception = new RuntimeException('exception caster regression');
$properties = Symfony\Component\VarDumper\Caster\Caster::castObject($exception, get_class($exception));
if (!is_int($properties["\0*\0line"])) {
    throw new RuntimeException('Exception line replaced by an uninitialized stub');
}
Symfony\Component\VarDumper\Caster\ExceptionCaster::castException(
    $exception,
    $properties,
    new Symfony\Component\VarDumper\Cloner\Stub(),
    false,
);
return true;
