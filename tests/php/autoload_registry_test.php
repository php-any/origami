<?php
class AutoloadRegistryProbe {
    public array $names = [];
    public function load(string $name): void { $this->names[] = $name; }
    public function __invoke(string $name): void { $this->names[] = 'invoke:'.$name; }
    public static function staticLoad(string $name): void {}
}
$first = new AutoloadRegistryProbe;
$second = new AutoloadRegistryProbe;
$closure = function ($name) {};
spl_autoload_register([$first, 'load']);
spl_autoload_register([$first, 'LOAD']);
spl_autoload_register($closure);
spl_autoload_register($second, true, true);
spl_autoload_register('AutoloadRegistryProbe::staticLoad');
spl_autoload_register(['AutoloadRegistryProbe', 'STATICLOAD']);
$callbacks = spl_autoload_functions();
if (count($callbacks) !== 4 || $callbacks[0] !== $second || $callbacks[1][0] !== $first || $callbacks[2] !== $closure) {
    throw new Exception('autoload identity, duplicate or prepend mismatch');
}
class_exists('AutoloadRegistryMissing');
if ($first->names !== ['AutoloadRegistryMissing'] || $second->names !== ['invoke:AutoloadRegistryMissing']) {
    throw new Exception('autoload callback scope mismatch');
}
if (!spl_autoload_unregister([$first, 'LOAD']) || spl_autoload_unregister([$first, 'load'])) {
    throw new Exception('autoload unregister mismatch');
}
if (!spl_autoload_unregister(['autoloadregistryprobe', 'staticload'])) { throw new Exception('static callback identity mismatch'); }
spl_autoload_unregister($closure);
spl_autoload_unregister($second);
if (count(spl_autoload_functions()) !== 0) { throw new Exception('autoload list retained removed callbacks'); }
spl_autoload_register([$first, 'load']);
return $first;
