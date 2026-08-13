<?php
namespace tests\php;

class NamedArg_App {
    public $last;
    public function register($provider, $force = false) {
        $this->last = [$provider, $force];
        return $provider;
    }
}

class NamedArg_Prov {}

$app = new NamedArg_App();
$path = __DIR__.'/../routes/web.php';
// simulate wrong capture? 
$app->register(NamedArg_Prov::class, force: true);
if ($app->last[0] !== NamedArg_Prov::class) {
    Log::fatal('provider wrong: '.var_export($app->last, true));
}
if ($app->last[1] !== true) {
    Log::fatal('force wrong: '.var_export($app->last, true));
}

// also test withRouting-like named args
function withRouting_like($using = null, $web = null, $api = null) {
    return [$using, $web, $api];
}
$web = '/tmp/web.php';
$r = withRouting_like(web: $web);
if ($r[1] !== $web || $r[0] !== null) {
    Log::fatal('withRouting named fail: '.var_export($r, true));
}

Log::info('named_arg_register 测试通过');
