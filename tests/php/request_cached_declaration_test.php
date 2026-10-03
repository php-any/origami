<?php
spl_autoload_register(function($name) {
    if ($name === 'CachedScopeTrait') { require __DIR__ . '/fixtures/cached_scope_trait.php'; }
});
class CachedScopeParent {
    public int $marker;
    public function __construct(int $marker) { $this->marker = $marker; }
}
class CachedScopeChild extends CachedScopeParent { use CachedScopeTrait; }
$marker = $request_marker ?? 41;
$object = new CachedScopeChild($marker);
CachedScopeChild::remember($marker);
if ($object->remembered() !== $marker || CachedScopeChild::$requestStack !== [$marker]) { throw new Exception('cached declaration state leaked'); }
return $object;
