<?php

namespace Illuminate\Foundation;

class AliasLoaderMini
{
    protected $aliases;
    protected $registered = false;
    protected static $facadeNamespace = 'Facades\\';
    protected static $instance;

    private function __construct($aliases)
    {
        $this->aliases = $aliases;
    }

    public static function getInstance(array $aliases = [])
    {
        if (is_null(static::$instance)) {
            return static::$instance = new static($aliases);
        }
        return static::$instance;
    }

    public function load($alias)
    {
        if (static::$facadeNamespace && str_starts_with($alias, static::$facadeNamespace)) {
            return true;
        }
        if (isset($this->aliases[$alias])) {
            return true;
        }
        return null;
    }

    public function register()
    {
        if (!$this->registered) {
            $this->prependToLoaderStack();
            $this->registered = true;
        }
    }

    protected function prependToLoaderStack()
    {
        spl_autoload_register([$this, 'load'], true, true);
    }
}

$loader = AliasLoaderMini::getInstance(['App' => 'Illuminate\\Support\\Facades\\App']);
$loader->register();

// direct call
$r = $loader->load('App');
echo "direct=".var_export($r, true)."\n";

// via class_exists / autoload
$ok = class_exists('Some\\MissingXYZ', true);
echo "autoload_ok\n";

// facade-ish
$r2 = $loader->load('Facades\\Something');
echo "facade=".var_export($r2, true)."\n";
