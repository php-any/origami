<?php

namespace Bootstrap\Console;

/**
 * 终端 ANSI 样式（支持 NO_COLOR 环境变量）
 */
class Style
{
    private static ?bool $colorsEnabled = null;

    public static function isColorEnabled(): bool
    {
        if (self::$colorsEnabled === null) {
            self::$colorsEnabled = getenv('NO_COLOR') === false;
        }

        return self::$colorsEnabled;
    }

    public static function wrap(string $text, string $code): string
    {
        if (!self::isColorEnabled()) {
            return $text;
        }

        return "\033[{$code}m{$text}\033[0m";
    }

    public static function bold(string $text): string
    {
        return self::wrap($text, '1');
    }

    public static function dim(string $text): string
    {
        return self::wrap($text, '2');
    }

    public static function green(string $text): string
    {
        return self::wrap($text, '32');
    }

    public static function yellow(string $text): string
    {
        return self::wrap($text, '33');
    }

    public static function red(string $text): string
    {
        return self::wrap($text, '31');
    }

    public static function cyan(string $text): string
    {
        return self::wrap($text, '36');
    }

    public static function blue(string $text): string
    {
        return self::wrap($text, '34');
    }
}
