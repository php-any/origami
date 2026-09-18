<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

// Same nested pattern via laravel13.exe
class NestedCaptureScope_Host2
{
    public function register()
    {
        return function (string $expression): string {
            [$name, $arguments] = str_contains($expression, ',') ?
                array_map('trim', explode(',', $expression, 2)) :
                [$expression, ''];
            return "<?php {$name} = X {$arguments} Y";
        };
    }
}

$fn = (new NestedCaptureScope_Host2())->register();
$out = $fn('$content, $logo, $isDarkMode = false');
echo "nested_out=".json_encode($out)."\n";
echo (str_contains($out, '$content =') ? "NESTED_OK\n" : "NESTED_BAD\n");

$compiler = app('blade.compiler');
$ref = new ReflectionProperty($compiler, 'customDirectives');
$ref->setAccessible(true);
$ffn = $ref->getValue($compiler)['capture'];
$fout = $ffn('$content, $logo, $isDarkMode = false');
echo "fil_out=".json_encode($fout)."\n";
echo (str_contains($fout, '$content =') ? "FIL_OK\n" : "FIL_BAD\n");
