<?php
class C {
    protected $routes = [];
    public function add() {
        $this->routes['GET']['/'] = 'hello';
    }
    public function getGet() {
        return $this->routes['GET'] ?? [];
    }
}
$c = new C();
$c->add();
$bucket = $c->getGet();
echo "bucket count=" . count($bucket) . "\n";
echo "val=" . ($bucket['/'] ?? 'missing') . "\n";
