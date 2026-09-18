<?php
/**
 * Subclass ComponentTagCompiler to dump attribute string and parsed attrs.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\View\Compilers\ComponentTagCompiler;

class DumpCTC extends ComponentTagCompiler
{
    public $dumps = [];

    protected function getAttributesFromAttributeString(string $attributeString)
    {
        $result = parent::getAttributesFromAttributeString($attributeString);
        $this->dumps[] = [
            'raw' => $attributeString,
            'raw_len' => strlen($attributeString),
            'raw_json' => json_encode($attributeString),
            'result' => $result,
            'bound' => $this->boundAttributes,
        ];
        return $result;
    }
}

$blade = app('blade.compiler');
$ctc = new DumpCTC(
    $blade->getClassComponentAliases(),
    $blade->getClassComponentNamespaces(),
    $blade
);

$value = '<x-filament-panels::page.simple>
    {{ $this->content }}
</x-filament-panels::page.simple>';

$out = $ctc->compile($value);
echo "DUMPS=".json_encode($ctc->dumps, JSON_UNESCAPED_UNICODE)."\n";
echo (str_contains($out, "'0'") ? "BAD\n" : "GOOD\n");
