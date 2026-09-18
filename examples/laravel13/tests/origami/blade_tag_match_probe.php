<?php
/**
 * Dump preg_replace_callback matches for Blade component opening tags.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$value = '<x-filament-panels::page.simple>
    {{ $this->content }}
</x-filament-panels::page.simple>';

$pattern = "/
            <
                \s*
                x[-\:]([\w\-\:\.]*)
                (?<attributes>
                    (?:
                        \s+
                        (?:
                            (?:
                                @(?:class)(\( (?: (?>[^()]+) | (?-1) )* \))
                            )
                            |
                            (?:
                                @(?:style)(\( (?: (?>[^()]+) | (?-1) )* \))
                            )
                            |
                            (?:
                                \{\{\s*\\\$attributes(?:[^}]+?)?\s*\}\}
                            )
                            |
                            (?:
                                (\:\\\$)(\w+)
                            )
                            |
                            (?:
                                [\w\-:.@%]+
                                (
                                    =
                                    (?:
                                        \\\"[^\\\"]*\\\"
                                        |
                                        \'[^\']*\'
                                        |
                                        [^\'\\\"=<>]+
                                    )
                                )?
                            )
                        )
                    )*
                    \s*
                )
                (?<![\/=\-])
            >
        /x";

if (!preg_match($pattern, $value, $matches)) {
    echo "NO_MATCH\n";
    exit;
}
echo "keys=".json_encode(array_keys($matches))."\n";
echo "1=".json_encode($matches[1] ?? null)."\n";
echo "attributes=".json_encode($matches['attributes'] ?? 'MISSING')."\n";
echo "attributes_len=".strlen((string)($matches['attributes'] ?? ''))."\n";
foreach ($matches as $k => $v) {
    if (is_string($k) || (is_int($k) && $k <= 5)) {
        echo "m[$k]=".json_encode($v)."\n";
    }
}
