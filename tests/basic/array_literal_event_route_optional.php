<?php

class ErRouteStub { public function getActionName(): string { return 'Home@index'; } }
class ErRequestStub { public function route() { return new ErRouteStub(); } }
class ErEventStub { public $request; public function __construct() { $this->request = new ErRequestStub(); } }

function er_optional($value)
{
    return new class ($value) {
        private $value;
        public function __construct($v) { $this->value = $v; }
        public function __call($m, $a) { return $this->value === null ? null : $this->value->$m(...$a); }
    };
}

class ErRW
{
    public function record($event): array
    {
        return [
            'controller_action' => er_optional($event->request->route())->getActionName(),
        ];
    }
}

$e = new ErEventStub();
$a = (new ErRW())->record($e);
if (($a['controller_action'] ?? '') !== 'Home@index') {
    Log::fatal('fail');
}
Log::info('event request route optional ok');
