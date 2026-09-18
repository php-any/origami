<?php
/**
 * Step through getAttributesFromAttributeString preprocessing on ''.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\View\Compilers\ComponentTagCompiler;

class StepCTC extends ComponentTagCompiler
{
    public function step(string $attributeString)
    {
        $steps = [];
        $steps[] = ['start', $attributeString, gettype($attributeString)];

        $pattern = "/\s\:\\\$(\w+)/x";
        $v = preg_replace_callback($pattern, function (array $matches) {
            return " :{$matches[1]}=\"\${$matches[1]}\"";
        }, $attributeString);
        $steps[] = ['short', $v, gettype($v), json_encode($v)];

        $pattern = "/
            (?:^|\s+)
            \{\{\s*(\\\$attributes(?:[^}]+?(?<!\s))?)\s*\}\}
        /x";
        $v2 = preg_replace($pattern, ' :attributes="$1"', $v);
        $steps[] = ['bag', $v2, gettype($v2), json_encode($v2)];

        $v3 = preg_replace_callback(
            '/@(class)(\( ( (?>[^()]+) | (?2) )* \))/x', function ($match) {
                if ($match[1] === 'class') {
                    $match[2] = str_replace('"', "'", $match[2]);
                    return ":class=\"\Illuminate\Support\Arr::toCssClasses{$match[2]}\"";
                }
                return $match[0];
            }, $v2
        );
        $steps[] = ['class', $v3, gettype($v3), json_encode($v3)];

        $v4 = preg_replace_callback(
            '/@(style)(\( ( (?>[^()]+) | (?2) )* \))/x', function ($match) {
                if ($match[1] === 'style') {
                    $match[2] = str_replace('"', "'", $match[2]);
                    return ":style=\"\Illuminate\Support\Arr::toCssStyles{$match[2]}\"";
                }
                return $match[0];
            }, $v3
        );
        $steps[] = ['style', $v4, gettype($v4), json_encode($v4)];

        $pattern = "/
            (?:^|\s+)
            :(?!:)
            ([\w\-:.@]+)
            =
        /xm";
        $v5 = preg_replace($pattern, ' bind:$1=', $v4);
        $steps[] = ['bind', $v5, gettype($v5), json_encode($v5)];

        $attrPattern = '/
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
        $n = preg_match_all($attrPattern, (string)$v5, $matches, PREG_SET_ORDER);
        $steps[] = ['match_n', $n, 'matches', $matches];

        return $steps;
    }
}

$blade = app('blade.compiler');
$ctc = new StepCTC($blade->getClassComponentAliases(), $blade->getClassComponentNamespaces(), $blade);
echo json_encode($ctc->step(''), JSON_PRETTY_PRINT), "\n";
