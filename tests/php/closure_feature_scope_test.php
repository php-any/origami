<?php
class ScopeEventBus {
    public static $listeners = [];
    public static function on($callback) { self::$listeners[] = $callback; }
    public function trigger($action) {
        $results = [];
        foreach (self::$listeners as $callback) { $results[] = $callback($action); }
        return $results;
    }
}
class ScopeMagicActions {
    public static $magicActions = ['$refresh', '$set'];
    public static function provide() {
        ScopeEventBus::on(function ($action) { return in_array($action, self::$magicActions); });
    }
}
class ScopeReleaseTokens {
    public static function provide() {
        ScopeEventBus::on(function ($action) { return self::class; });
    }
}
class ScopeMagicChild extends ScopeMagicActions { public static $magicActions = ['child']; }
foreach ([ScopeMagicActions::class, ScopeReleaseTokens::class, ScopeMagicChild::class] as $feature) { $feature::provide(); }
$bus = new ScopeEventBus();
$result = $bus->trigger('$refresh');
if ($result !== [true, ScopeReleaseTokens::class, true]) { throw new Exception('feature callback lexical scope'); }
echo "closure feature scope: PASS\n";
