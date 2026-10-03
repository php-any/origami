<?php
trait CachedScopeTrait {
    public static array $requestStack = [];
    public function remembered(): int { return $this->marker; }
    public static function remember(int $value): void { array_push(self::$requestStack,$value); }
}
