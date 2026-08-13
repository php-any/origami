<?php

class AfCollection
{
    private $x;
    public function __construct($x) { $this->x = $x; }
    public function map($cb) {
        $out = [];
        foreach ($this->x as $k => $v) { $out[$k] = $cb($v); }
        return new AfCollection($out);
    }
    public function all() { return $this->x; }
}

function af_collect($x) {
    return new AfCollection($x);
}

$methods = ['GET', 'POST'];
$result = af_collect($methods)->map(fn ($method) => strtolower($method))->all();
if ($result !== ['get', 'post']) {
    Log::fatal('arrow map fail');
}
Log::info('arrow fn in collect map ok');
