<?php
register_shutdown_function(function () {
    echo "first\n";
    register_shutdown_function(function () { echo "late\n"; });
});
register_shutdown_function(function () { echo "second\n"; });
function shutdown_named($value) { echo $value."\n"; }
register_shutdown_function('shutdown_named', 'named');
class ShutdownReceiver {
    public function __invoke($value) { echo $value."\n"; }
}
register_shutdown_function(new ShutdownReceiver, 'invokable');
echo "body\n";
