<?php

class OcRouteStub
{
    public function getActionName(): string
    {
        return 'App\\Http\\Controllers\\HomeController@index';
    }
}

function oc_optional($value)
{
    return new class ($value) {
        private $value;

        public function __construct($value)
        {
            $this->value = $value;
        }

        public function __call($method, $args)
        {
            if ($this->value === null) {
                return null;
            }
            return $this->value->$method(...$args);
        }
    };
}

class OcRW
{
    public function record(): array
    {
        $route = new OcRouteStub();

        return [
            'controller_action' => oc_optional($route)->getActionName(),
            'method' => 'GET',
        ];
    }
}

$w = new OcRW();
$a = $w->record();
if (($a['controller_action'] ?? '') !== 'App\\Http\\Controllers\\HomeController@index') {
    Log::fatal('optional chain in array literal failed');
}
if (($a['method'] ?? '') !== 'GET') {
    Log::fatal('second key failed');
}
Log::info('array literal optional chain ok');
