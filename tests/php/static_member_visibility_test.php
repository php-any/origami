<?php
class StaticVisibilityParent {
    protected static mixed $instance = null;
    private static int $secret = 4;
    public static function getInstance(): mixed { return static::$instance ??= new static; }
    public static function getSecret(): int { return self::$secret; }
}
class StaticVisibilityChild extends StaticVisibilityParent {}
if (!(StaticVisibilityChild::getInstance() instanceof StaticVisibilityChild) || StaticVisibilityChild::getSecret() !== 4) { throw new Exception('static lexical visibility'); }
try { $x = StaticVisibilityParent::$secret; throw new Exception('static privacy lost'); } catch (Error $e) {}
echo "static member visibility OK\n";
