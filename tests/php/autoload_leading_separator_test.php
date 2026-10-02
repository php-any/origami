<?php

$loader = function ($class) {
    if ($class === 'AutoloadLeadingSeparator\Admin') {
        require __DIR__.'/autoload_leading_separator_fixtures/Admin.php';
    } elseif ($class === 'AutoloadLeadingSeparator\Contract') {
        require __DIR__.'/autoload_leading_separator_fixtures/Contract.php';
    } else {
        throw new \Exception('Unexpected autoload name: '.$class);
    }
};
spl_autoload_register($loader);

// EloquentUserProvider::createModel uses a leading separator for dynamic new.
$class = '\AutoloadLeadingSeparator\Admin';
$admin = new $class;
if (get_class($admin) !== 'AutoloadLeadingSeparator\Admin') {
    throw new \Exception('Dynamic new resolved the wrong class');
}
if (!interface_exists('\AutoloadLeadingSeparator\Contract')) {
    throw new \Exception('Interface autoload failed');
}

spl_autoload_unregister($loader);
echo "autoload leading separator OK\n";
