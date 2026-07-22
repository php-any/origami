<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Bus\Dispatcher;
use Illuminate\Container\Container;

class BusSmoke_Job
{
    public $value;

    public function __construct($value)
    {
        $this->value = $value;
    }

    public function handle()
    {
        return 'handled:' . $this->value;
    }
}

$container = new Container();
$bus = new Dispatcher($container);

$result = $bus->dispatchNow(new BusSmoke_Job('origami'));
if ($result !== 'handled:origami') {
    echo "FAIL: dispatchNow got ";
    var_export($result);
    echo "\n";
    exit(1);
}

echo "PASS\n";
