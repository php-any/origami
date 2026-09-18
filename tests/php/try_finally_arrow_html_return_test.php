<?php
/**
 * 最小复现：try/finally + 箭头 + 声明 string 返回。
 */
class Component {
    public static bool $cache = false;
    public static function withVisibilityCache(callable $callback): mixed {
        $wasEnabled = self::$cache;
        if (!$wasEnabled) self::$cache = true;
        try {
            return $callback();
        } finally {
            if (!$wasEnabled) self::$cache = false;
        }
    }
}

class Schema {
    public function toEmbeddedHtml(): string {
        return Component::withVisibilityCache(fn (): string => $this->renderEmbeddedHtml());
    }
    protected function renderEmbeddedHtml(): string {
        return 'ok';
    }
}

class SchemaHtml {
    public function toEmbeddedHtml(): string {
        return Component::withVisibilityCache(fn (): string => $this->renderEmbeddedHtml());
    }
    protected function renderEmbeddedHtml(): string {
        if (false) return '';
        ob_start(); ?>
        <div>hi</div>
        <?php return ob_get_clean();
    }
}

$s = new Schema();
$r = $s->toEmbeddedHtml();
echo "plain=".gettype($r)." val=$r\n";

$h = new SchemaHtml();
try {
    $r = $h->toEmbeddedHtml();
    echo "html=".gettype($r)." val=".var_export($r,true)."\n";
    if (is_object($r)) echo "html_class=".get_class($r)."\n";
} catch (Throwable $e) {
    echo "html EX: ".$e->getMessage()."\n";
}

// 直接 render
$ref = new ReflectionClass($h);
$m = $ref->getMethod('renderEmbeddedHtml');
$m->setAccessible(true);
try {
    $r = $m->invoke($h);
    echo "direct_render=".gettype($r)." val=".var_export($r,true)."\n";
} catch (Throwable $e) {
    echo "direct EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
