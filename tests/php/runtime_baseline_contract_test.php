<?php
function baseline_check($ok, $message) { if (!$ok) throw new Exception($message); }
baseline_check(dirname('/a//b/c/', 2) === '/a', 'dirname preserves input separators');
baseline_check(dirname('/a/../b') === '/a/..', 'dirname is lexical');
try { dirname('x', 0); throw new Exception('invalid levels accepted'); } catch (ValueError $e) {}
function baseline_join(string $a, string $b): string { return $a . ':' . $b; }
baseline_check(baseline_join(...['b'=>'Y', 'a'=>'X']) === 'X:Y', 'named unpack');
try { baseline_join('X', ...['a'=>'Y']); throw new Exception('duplicate accepted'); } catch (Error $e) {}
try { match (0) { 1 => 'one' }; throw new Exception('unmatched accepted'); } catch (UnhandledMatchError $e) {}
baseline_check(0o17 === 15 && 0xFE === 254, 'number bases');
$r = new ReflectionFunction('baseline_join');
baseline_check($r->getParameters()[0]->getType()->getName() === 'string', 'named function reflection');
$warnings = [];
set_error_handler(function($code,$message) use (&$warnings) { $warnings[]=$code; return true; });
class BaselinePointer implements Iterator {
    private $first = 10; public $second = 20;
    function current(): mixed { throw new Exception('iterator called'); }
    function key(): mixed { throw new Exception('iterator called'); }
    function next(): void { throw new Exception('iterator called'); }
    function rewind(): void { throw new Exception('iterator called'); }
    function valid(): bool { throw new Exception('iterator called'); }
}
$o = new BaselinePointer();
baseline_check(reset($o) === 10 && key($o) === "\0BaselinePointer\0first", 'object property pointer');
baseline_check(next($o) === 20 && end($o) === 20 && prev($o) === 10, 'object pointer movement');
baseline_check($warnings === [8192,8192,8192,8192,8192], 'object pointer deprecations');
restore_error_handler();
$file = __DIR__.'/../../examples/laravel13/storage/origami-debug/runtime-baseline-source.php';
file_put_contents($file, '<?php return 1;');
baseline_check((require $file) === 1, 'first include');
file_put_contents($file, '<?php return 2;');
baseline_check((require $file) === 2, 'rewritten include');
$stream = fopen($file, 'w');
fwrite($stream, '<?php return 3;');
fclose($stream);
baseline_check((require $file) === 3, 'stream rewritten include');
$replacement = $file.'.replacement';
file_put_contents($replacement, '<?php return 4;');
rename($replacement, $file);
baseline_check((require $file) === 4, 'renamed include');
unlink($file);
echo "runtime baseline contract OK\n";
