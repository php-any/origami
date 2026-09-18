<?php
namespace tests\php;

class Comp {
    public $content = 'from-comp';
}
class FS {
    public function getRequire($path) {
        return (static function () use ($path) {
            // In PHP, $this is unavailable in static closures
            try {
                return isset($this);
            } catch (\Error $e) {
                return 'no-this:'.$e->getMessage();
            }
        })();
    }
    public function bindRequire($path, $comp) {
        return \Closure::bind(function () use ($path) {
            return $this->content;
        }, $comp, $comp)();
    }
}

$fs = new FS();
$r = $fs->getRequire('x');
if ($r === true || $r === false) {
    // Origami might return bool for isset($this)
    Log::info('static isset this='.var_export($r, true));
} else {
    Log::info('static: '.$r);
}

$c = new Comp();
$got = $fs->bindRequire('x', $c);
if ($got !== 'from-comp') {
    Log::fatal('Closure::bind $this failed: '.var_export($got, true));
}
Log::info('Closure::bind ok');
