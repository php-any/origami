<?php
/**
 * Inline copy of getAttributesFromAttributeString with dumps.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;
use Illuminate\Support\Str;
use Illuminate\View\Compilers\ComponentTagCompiler;

class InlineCTC extends ComponentTagCompiler
{
    public function debugAttrs(string $attributeString)
    {
        $log = [];
        $attributeString = $this->parseShortAttributeSyntax($attributeString);
        $log[] = ['after_short', $attributeString, gettype($attributeString)];
        $attributeString = $this->parseAttributeBag($attributeString);
        $log[] = ['after_bag', $attributeString, gettype($attributeString)];
        $attributeString = $this->parseComponentTagClassStatements($attributeString);
        $log[] = ['after_class', $attributeString, gettype($attributeString)];
        $attributeString = $this->parseComponentTagStyleStatements($attributeString);
        $log[] = ['after_style', $attributeString, gettype($attributeString)];
        $attributeString = $this->parseBindAttributes($attributeString);
        $log[] = ['after_bind', $attributeString, gettype($attributeString)];

        $pattern = '/
            (?<attribute>[\w\-:.@%]+)
            (
                =
                (?<value>
                    (
                        \"[^\"]+\"
                        |
                        \\\'[^\\\']+\\\'
                        |
                        [^\s>]+
                    )
                )
            )?
        /x';

        $matches = [];
        $n = preg_match_all($pattern, $attributeString, $matches, PREG_SET_ORDER);
        $log[] = ['n' => $n, 'matches' => $matches, 'n_type' => gettype($n)];

        if (! $n) {
            return [$log, []];
        }

        $bound = [];
        $result = (new Collection($matches))->mapWithKeys(function ($match) use (&$bound) {
            $attribute = $match['attribute'];
            $value = $match['value'] ?? null;
            if (is_null($value)) {
                $value = 'true';
                $attribute = Str::start($attribute, 'bind:');
            }
            $value = $this->stripQuotes($value);
            if (str_starts_with($attribute, 'bind:')) {
                $attribute = Str::after($attribute, 'bind:');
                $bound[$attribute] = true;
            } else {
                $value = "'".$this->compileAttributeEchos($value)."'";
            }
            if (str_starts_with($attribute, '::')) {
                $attribute = substr($attribute, 1);
            }
            return [$attribute => $value];
        })->toArray();

        $log[] = ['result' => $result, 'bound' => $bound];
        return [$log, $result];
    }
}

$blade = app('blade.compiler');
$ctc = new InlineCTC($blade->getClassComponentAliases(), $blade->getClassComponentNamespaces(), $blade);
[$log, $result] = $ctc->debugAttrs('');
echo json_encode($log, JSON_PRETTY_PRINT|JSON_UNESCAPED_UNICODE), "\n";
echo "RESULT=".json_encode($result)."\n";

// Compare to real parent
$ref = new ReflectionClass(ComponentTagCompiler::class);
$m = $ref->getMethod('getAttributesFromAttributeString');
$m->setAccessible(true);
$real = $m->invoke($ctc, '');
echo "REAL=".json_encode($real)."\n";
echo "BOUND=".json_encode($ctc->boundAttributes ?? null)."\n";
