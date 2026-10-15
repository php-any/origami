<?php
if (__LINE__ !== 2) { throw new Exception('top-level __LINE__ must be one-based'); }
function magic_line_value() {
    return __LINE__;
}
if (magic_line_value() !== 4) { throw new Exception('function __LINE__ must retain its source line'); }
trait MagicLineTrait {
    public function sourceLine() { return __LINE__; }
}
class MagicLineClass { use MagicLineTrait; }
if ((new MagicLineClass)->sourceLine() !== 8) { throw new Exception('trait method __LINE__ must retain its declaration line'); }
echo "magic line one-based: PASS\n";
