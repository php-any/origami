<?php

require __DIR__.'/../../vendor/autoload.php';

class UtilsProbe2 {
    static function escapeStringForHtml($subject) {
        return 'E:'.(string)$subject;
    }
    static function stringifyHtmlAttributes($attributes) {
        return collect($attributes)
            ->mapWithKeys(function ($value, $key) {
                return [$key => static::escapeStringForHtml($value)];
            })->map(function ($value, $key) {
                return sprintf('%s="%s"', $key, $value);
            })->implode(' ');
    }
}

try {
    echo UtilsProbe2::stringifyHtmlAttributes(['id' => 'x'])."\n";
} catch (Throwable $e) {
    echo "ERR ".$e->getMessage()." at ".$e->getFile().":".$e->getLine()."\n";
}
