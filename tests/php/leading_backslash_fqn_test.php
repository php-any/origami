<?php
namespace tests\php;
class WireKeysProbe {
    public static function openLoop() { return 'ok'; }
}
// FQN with leading backslash like compiled blade
echo \tests\php\WireKeysProbe::openLoop(), "\n";
