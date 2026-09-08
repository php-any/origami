<?php

use Illuminate\Support\Collection;
use Egulias\EmailValidator\Validation\RFCValidation;
use Egulias\EmailValidator\Validation\MultipleValidationWithAnd;
use Egulias\EmailValidator\EmailValidator;
use Illuminate\Container\Container;

require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

$parameters = [];
$validations = (new Collection($parameters))
    ->unique()
    ->map(fn ($validation) => match (true) {
        $validation === 'strict' => 'strict',
        default => new RFCValidation(),
    })
    ->all() ?: [new RFCValidation];

echo "count=".count($validations)." type=".gettype($validations)."\n";
echo "is_array=".var_export(is_array($validations), true)."\n";
echo "empty_elvis=".var_export(([] ?: ['x']), true)."\n";

$all = (new Collection($parameters))->unique()->map(fn ($v) => new RFCValidation())->all();
echo "all_count=".count($all)." all_bool=".var_export((bool)$all, true)."\n";
$after = $all ?: [new RFCValidation];
echo "after_count=".count($after)."\n";

try {
    $mv = new MultipleValidationWithAnd($after);
    echo "MultipleValidationWithAnd=ok\n";
} catch (Throwable $e) {
    echo "ERR ".$e->getMessage()."\n";
}

// What does Validator email rule pass as parameters?
$v = validator(['email' => 'admin@example.com'], ['email' => 'email']);
try {
    echo "passes=".var_export($v->passes(), true)."\n";
    echo "errors=".json_encode($v->errors()->all())."\n";
} catch (Throwable $e) {
    echo "validator_err=".get_class($e)." ".$e->getMessage()."\n";
}
