<?php
class VisibilityParent {
    private int $secret = 7;
    protected int $shared = 8;
    private static int $hidden = 9;
    private function secret(): string { return 'private'; }
    protected function shared(): string { return 'protected'; }
    public function own(): array { return [$this->secret(), $this->secret, self::$hidden]; }
    public function privateClosure(): Closure { return $this->secret(...); }
}
class VisibilityChild extends VisibilityParent {
    public function parentPrivate(): mixed { return $this->secret(); }
    public function parentPrivateProperty(): mixed { return $this->secret; }
    public function parentProtected(): array { return [$this->shared(), $this->shared]; }
    public function parentStatic(): mixed { return self::$hidden; }
}
class VisibilityMagic extends VisibilityParent {
    public array $writes = [];
    public function __call(string $name, array $args): string { return 'magic:' . $name; }
    public function __get(string $name): mixed { return 'get:' . $name; }
    public function __set(string $name, mixed $value): void { $this->writes[$name] = $value; }
    public function fromChild(): array { return [$this->secret(), $this->secret]; }
}
class VisibilityShadow extends VisibilityParent {
    public int $secret = 30;
    public function secret(): string { return 'child'; }
}
$shadow = new VisibilityShadow();
if ($shadow->own() !== ['private', 7, 9] || $shadow->secret !== 30) { throw new Exception('private member collision'); }
$child = new VisibilityChild();
if ($child->own() !== ['private', 7, 9] || $child->parentProtected() !== ['protected', 8] || $child->privateClosure()() !== 'private') { throw new Exception('lexical access failed'); }
$errors = 0;
foreach (['parentPrivate', 'parentStatic'] as $name) {
    try { $child->$name(); } catch (Error $e) { $errors++; }
}
try { $invalid = $child->secret(...); } catch (Error $e) { $errors++; }
$parent = new VisibilityParent();
try { $parent->secret = 'wrong type'; } catch (Error $e) { $errors++; }
if ($errors !== 4) { throw new Exception('private access escaped'); }
$magic = new VisibilityMagic();
if ($magic->secret() !== 'magic:secret' || $magic->secret(...)() !== 'magic:secret' || $magic->fromChild() !== ['magic:secret', 'get:secret']) { throw new Exception('magic visibility failed'); }
$name = 'secret';
$magic->$name = 'different type';
if ($magic->$name !== 'get:secret' || $magic->writes[$name] !== 'different type') { throw new Exception('dynamic visibility failed'); }
echo "member visibility matrix OK\n";
